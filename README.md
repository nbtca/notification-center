# webhook-delivery-center

webhook 转发中心

- 在配置的端口监听 http(s)请求，用于接收 webhook 请求
- 在配置的端口/ws 路径下监听 websocket 请求，连接到此 websocket 的客户端会收到转发自 webhook 的请求消息

## 🚀 Deployment/部署

- ### 🐳 Docker

  执行`docker_build.sh`，将构建并打包`webhook.tar`镜像文件，在运行环境执行`docker load -i webhook.tar`导入镜像

  - 启动命令(参考)
    ```bash
    touch webhook-delivery-center.config.json
    docker run -d -p 8080:8080 --name webhook -v $(pwd)/webhook-delivery-center.config.json:/config/config.json webhook
    ```

## 🛠️ Config/配置

```json
{
  "bind": ":8080",
  "use_cert": false,
  "cert_file": "fullchain.cer",
  "key_file": "private.key",
  "auth": {
    "": "默认路径的密钥",
    "github":"对于/github路径的请求使用的密钥"
  }
}
```

## 🛡️ Authentication/鉴权

> 不论是 http(s) 还是 websocket 请求，都需要进行鉴权，鉴权方式如下二选一

- ### Header["Authorization"]
  - 用于直接鉴权的密钥
  - 格式为`Bearer ${value}`，`${value}`为配置文件中`auth`字段中的`value`对应的值
- ### Header["X-Signature-256"]
  - `SHA256`签名，用于验证请求的合法性，通过`auth`字段中的`value`作为 HMAC 的密钥，对请求的`body`进行签名，签名结果与请求头中的`X-Signature-256`进行比对，如果一致则请求合法，否则请求非法
