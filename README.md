# github-webhook-proxy

simple GitHub webhook proxy

此程序运行于公网服务器

- 在配置的端口监听 http(s)请求，用于接收 GitHub 的 webhook 请求，
- 在配置的端口/ws 路径下监听 websocket 请求，连接到此 websocket 的(内网)客户端会收到转发自 GitHub Webhook 的请求消息
