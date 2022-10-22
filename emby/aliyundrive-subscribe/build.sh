#!/bin/bash

chmod 777 ./aliyundrive-subscribe_linux_amd64
docker buildx build --platform=linux/amd64 -t yqchilde/aliyundrive-subscribe .
docker tag yqchilde/aliyundrive-subscribe yqchilde/aliyundrive-subscribe:linux-amd64
docker push yqchilde/aliyundrive-subscribe:linux-amd64
