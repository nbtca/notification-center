# webhook-delivery-center

webhook 转发中心

- 在配置的端口监听 http(s)请求，用于接收 webhook 请求
- 在配置的端口/ws 路径下监听 websocket 请求，连接到此 websocket 的客户端会收到转发自 webhook 的请求消息
