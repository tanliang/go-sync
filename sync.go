package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var syncDir string

func main() {

	syncDir, _ = os.Getwd()
	syncDir += string(os.PathSeparator) + "tmp" + string(os.PathSeparator)
	if !fileExist(syncDir) {
		os.Mkdir(syncDir, os.ModePerm)
	}

	// clean expired item
	ticker := time.NewTicker(time.Second * 60 * 60)
	go func() {
		for {
			<-ticker.C
			walkDir(syncDir, 60 * 15)
		}
	}()

	// Echo instance
	e := echo.New()

	log.Print("=============sync_dir============")
	log.Print(syncDir)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://9777.in"},
		//AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	e.POST("/go/sync", sync)
	e.GET("/go/sync", sync)
	e.POST("/go/test", test)
	e.POST("/go", home)
	e.GET("/go", home)

	// Start server
	e.Logger.Fatal(e.Start(":9501"))
}

type ReqSync struct {
	Id string `query:"id" form:"id" json:"id"`
	Text string `query:"text" form:"text" json:"text"`
}


// Handler
func home(c echo.Context) error {
	rand.Seed(time.Now().Unix())
	id := strconv.FormatInt(int64(rand.Uint32()), 16)
	resp := map[string]string{
		"code": "200",
		"id": id,
	}
	return c.JSON(http.StatusOK, resp)
}

// Handler
func test(c echo.Context) error {
	reqSync := new(ReqSync)
	if err := c.Bind(reqSync); err != nil {
		return err
	}
	success := ""
	if reqSync.Id != "" {
		id, _ := strconv.ParseInt(reqSync.Id, 16, 32)
		fp := syncDir +strconv.FormatInt(id, 10)+".txt"
		if fileExist(fp) {
			success = "1"
		}
	}

	resp := map[string]string{
		"code": "200",
		"success": success,
	}
	return c.JSON(http.StatusOK, resp)
}

func sync(c echo.Context) error {
	reqSync := new(ReqSync)
	if err := c.Bind(reqSync); err != nil {
		return err
	}

	rspText := ""
	if reqSync.Id != "" {
		id, _ := strconv.ParseInt(reqSync.Id, 16, 32)
		fp := syncDir +strconv.FormatInt(id, 10)+".txt"

		if reqSync.Text != "" {
			content := []byte(reqSync.Text)
			err := ioutil.WriteFile(fp, content, 0644)
			if err != nil {
				panic(err)
			}
			rspText = reqSync.Text
		} else {
			file, err := os.Open(fp)
			if err == nil {
				//panic(err)
				content, _ := ioutil.ReadAll(file)
				rspText = string(content)
			}
			defer file.Close()
		}
	}

	resp := map[string]string{
		"code": "200",
		"text": rspText,
	}
	return c.JSON(http.StatusOK, resp)
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func walkDir(dir string, suffix int64) (err error) {
	nowTime := time.Now().Unix()
	err = filepath.Walk(dir, func(fname string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		fileTime := fi.ModTime().Unix()

		if (nowTime - fileTime) > suffix {
			remove(fname)
		}

		return nil
	})

	return err
}

func remove(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}