
## 单体应用分支
master

## 分布式分支

1. udp 进行数据分发的分支 dis_master
2. nginx 进行反向代理
1 2 配合
```

#user  nobody;
worker_processes  1;

#error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

#pid        logs/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    
    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;
	upstream wsbackend {
			server 192.168.0.102:8080;
			server 192.168.0.100:8080;
			hash $request_uri;
	}
	map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
	}
    server {
	  listen  80;
	  server_name localhost;
	  location / {
	   proxy_pass http://wsbackend;
	  }
	  #
	  location ^~ /chat {
	   proxy_pass http://wsbackend;
	   #防止连接中断
	   proxy_connect_timeout 500s;
       proxy_read_timeout 500s;
	   proxy_send_timeout 500s;
	   
	   ## 需要支持 websocket,设置一下两个参数
	   proxy_set_header Upgrade $http_upgrade;	#websockt 的字符串
       proxy_set_header Connection "Upgrade";
	  }
	 }

}

```

## 打包发布
build.bat 打包 window 平台
build.sh 打包 mac 平台

注意在 linux 下执行
`nohup ./chat >>./log.log 2>&1 &`