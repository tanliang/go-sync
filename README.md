# go-sync
a tiny web service for sending message via qr-code, mainly used in tv box.

# install
put index.html and sync.html under nginx root dir

# nginx
set CORS and Reverse Proxy
```
server {
    ...
    
    location /go {
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
        add_header Access-Control-Allow-Headers 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization';

        proxy_pass http://127.0.0.1:9501;

        if ($request_method = 'OPTIONS') {
            return 204;
        }
    }
    
    ...
}
```

# example
http://9777.in

# 说明
手机与电视盒子之间同步文本的服务