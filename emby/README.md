## 由来

方案来自群佬，他的仓库地址 [https://github.com/zxcvbnmzsedr/docker_env/blob/master/emby/README.md](https://github.com/zxcvbnmzsedr/docker_env/blob/master/emby/README.md)

## 挂载webdav

使用 `rclone` 工具连接webdav，使用 `systemctl` 做开机自启，`rclone.service` 内容如下：

```text
[Unit]
Description=rclone
After=docker.service

[Service]
User=root
# 这是注释，写入时删掉，`aliyun:/` 这一部分要按需变更，准确的话要改成rclone挂载盘的具体路径
ExecStart=/usr/bin/rclone mount aliyun:/ /mnt/aliyundrive --allow-other --allow-non-empty --vfs-cache-mode writes
Restart=on-abort

[Install]
WantedBy=multi-user.target
```

## 注意

1. `docker-compose.yaml` 中挂载项 `- /mnt/aliyundrive:/media` ，要修改为自己的云盘挂载目录
2. 需要修改 `nginx/conf.d/emby.js` 文件来配置云盘直链获取
3. `Emby.Plugins.Douban.dll` 文件为豆瓣刮削插件，可导入 `emby/config/plugins` 目录下
4. `docker-compose.yaml` 容器中 `tmm` 服务配置项 `extra_hosts` 需要改为自己最快的 `tvdb` dns地址
    * 查询地址可在 [https://dnschecker.org](https://dnschecker.org) 查询