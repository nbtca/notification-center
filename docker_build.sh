# build image
docker build -t notification-center .
# export image
docker save -o notification-center.tar notification-center
# import image
# docker load -i notification-center.tar
