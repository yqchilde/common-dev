#!/bin/bash

docker rm -f aliyunsub

docker run -d \
    --name aliyunsub \
    -p 8002:8002 \
    -v "$(pwd)"/conf:/app/conf \
    -e TZ=Asia/Shanghai \
    --restart unless-stopped \
    yqchilde/aliyundrive-subscribe:linux-amd64
