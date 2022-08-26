#!/bin/bash

# ./run.sh start 开启代理
# ./run.sh stop 关闭代理
function githubProxy() {
    if [[ "$1" = "start" ]]; then
        sudo /usr/bin/fastGithub start
        export http_proxy=http://127.0.0.1:38457
        export https_proxy=http://127.0.0.1:38457
        export socks_proxy=http://127.0.0.1:38457
        echo "github proxy 代理开启，端口：38457"
    elif [[ "$1" = "stop" ]]; then
        sudo /usr/bin/fastGithub stop
        unset http_proxy
        unset https_proxy
        unset socks_proxy
        echo "github proxy 代理关闭"
    fi
}


