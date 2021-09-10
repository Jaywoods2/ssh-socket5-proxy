```graphLR
    A[Local] -->|Ssh Tunnel| B(Middle Server:22)
    B --> C[Pod1]
    B --> D[Pod2]
```

测试:
```bash
curl --proxy "socks5://localhost:10111" 10.244.61.35:80
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

```bigquery
2021/09/10 18:49:20 target: 10.244.61.35:80
2021/09/10 18:49:20 192.168.33.1:63816  ===>  192.168.33.14:22 连接建立成功
2021/09/10 18:49:20 192.168.33.14:22  ===>  10.244.61.35:80 连接建立成功
2021/09/10 18:49:20 入流量共76.00B
2021/09/10 18:50:14 出流量共850.00B
```