FROM nginx
MAINTAINER yqchilde

WORKDIR /app
ADD ./conf/ /app/conf/
ADD ./aliyundrive-subscribe_linux_amd64 /app/aliyundrive-subscribe_linux_amd64

EXPOSE 8002
ENTRYPOINT ["./aliyundrive-subscribe_linux_amd64"]