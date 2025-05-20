# Nginx日志目录

此目录用于存放Nginx的访问日志和错误日志文件。

## 日志文件说明

- `access.log`: 记录所有对服务器的请求
- `error.log`: 记录Nginx运行时的错误信息

## 日志配置

在Nginx配置文件中，日志通常配置如下：

```nginx
server {
    # 其他配置...
    
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;
    
    # 其他配置...
}
```

这些日志文件对于监控服务器性能、排查问题和安全审计非常重要。