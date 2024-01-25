# check if webhook_test is running
if [ "$(docker ps -q -f name=webhook_test)" ]; then
    docker stop webhook_test
fi
# remove webhook_test container if it exists
if [ "$(docker ps -aq -f status=exited -f name=webhook_test)" ]; then
    docker rm webhook_test
fi
# 将本地的 'webhook-delivery-center.config.json' 映射到容器内的 '/config/config.json'
docker run -d -p 18080:8080 --name webhook_test -v $(pwd)/webhook-delivery-center.config.json:/config/config.json webhook

