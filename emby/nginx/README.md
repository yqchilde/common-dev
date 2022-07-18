## 由来

方案来自群佬，他的仓库地址 [https://github.com/zxcvbnmzsedr/docker_env/blob/master/emby/README.md](https://github.com/zxcvbnmzsedr/docker_env/blob/master/emby/README.md)

## 注意

1. `docker-compose.yaml` 中挂载项 `- /mnt/aliyundrive:/media` ，要修改为自己的云盘挂载目录
2. 需要修改 `nginx/conf.d/emby.js` 文件来配置云盘直链获取
3. `Emby.Plugins.Douban.dll` 文件为豆瓣刮削插件，可导入 `emby/config/plugins` 目录下
4. `docker-compose.yaml` 容器中 `tmm` 服务配置项 `extra_hosts` 需要改为自己最快的 `tvdb` dns地址
    * 查询地址可在 [https://dnschecker.org](https://dnschecker.org) 查询