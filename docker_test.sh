# check if notification-center_test is running
if [ "$(docker ps -q -f name=notification-center_test)" ]; then
    docker stop notification-center_test
fi
# remove notification-center_test container if it exists
if [ "$(docker ps -aq -f status=exited -f name=notification-center_test)" ]; then
    docker rm notification-center_test
fi
# 将本地的 'notification-center.config.json' 映射到容器内的 '/config/config.json'
docker run -d -p 18080:8080 --name notification-center_test -v $(pwd)/notification-center.config.json:/config/config.json notification-center

