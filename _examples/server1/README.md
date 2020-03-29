# 演示实例1

该示例支持http 协议和 TCP自定义协议

1.运行：
```
go run main.go
```

2.浏览器访问：
http://127.0.0.1:8090/a?hello   
看到内容：
```
你好:/a?hello
```

3.telnet访问自定义协议：
```
telnet 127.0.0.1 8080
```

输入：
```
say:hello
```
服务端响应：
```
reply:a
```