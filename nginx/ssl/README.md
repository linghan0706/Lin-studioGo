# SSL证书目录

此目录用于存放SSL证书文件。

## 使用说明

1. 将您的SSL证书文件（通常为.crt或.pem格式）放在此目录中
2. 将您的SSL私钥文件（通常为.key格式）放在此目录中
3. 在Nginx配置文件中引用这些证书文件

## 示例配置

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/your-certificate.crt;
    ssl_certificate_key /etc/nginx/ssl/your-private-key.key;
    
    # 其他SSL配置...
}
```

注意：在生产环境中，请确保您的SSL证书是由受信任的证书颁发机构签发的。