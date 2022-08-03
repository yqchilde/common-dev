// 查看日志: docker logs --tail 100 -f embynginx 2>&1 | grep "js:"

async function redirect2Pan(r) {
    // 这里默认emby/jellyfin的地址是宿主机
    const embyHost = 'http://192.168.31.7:8096';

    // rclone 的挂载目录, 在docker-compose中已将rclone挂载到/media目录下，故这里写/media
    const embyMountPath = '/media';

    // alist password
    const alistPwd = 'xxx';

    // 访问宿主机上5244端口的alist api
    const alistApiPath = 'http://192.168.31.5:5244/api/public/path';

    // 指定alist中目录名，这样会指定扫描该目录，若不写会扫描整个目录
    const alistPath = "AliYun"

    // fetch mount emby/jellyfin file path
    const regex = /[A-Za-z0-9]+/g;
    const itemId = r.uri.replace('emby', '').replace(/-/g, '').match(regex)[1];
    const mediaSourceId = r.args.MediaSourceId;
    let api_key = r.args.api_key;

    // infuse用户需要填写下面的api_key, 感谢@amwamw968
    if ((api_key === null) || (api_key === undefined)) {
        // 这里填自己的emby/jellyfin API KEY
        api_key = 'xxx';
        r.warn(`api key for Infuse: ${api_key}`);
    }

    const itemInfoUri = `${embyHost}/Items/${itemId}/PlaybackInfo?MediaSourceId=${mediaSourceId}&api_key=${api_key}`;
    r.warn(`itemInfoUri: ${itemInfoUri}`);
    const embyRes = await fetchEmbyFilePath(itemInfoUri);
    if (embyRes.startsWith('error')) {
        r.error(embyRes);
        r.return(500, embyRes);
        return;
    }
    r.warn(`mount emby file path: ${embyRes}`);

    // fetch alist direct link
    const alistFilePath = embyRes.replace(embyMountPath, alistPath);
    const alistRes = await fetchAlistPathApi(alistApiPath, alistFilePath, alistPwd);
    if (!alistRes.startsWith('error')) {
        r.warn(`redirect to: ${alistRes}`);
        r.return(302, alistRes);
        return;
    }
    if (alistRes.startsWith('error401')) {
        r.error(alistRes);
        r.return(401, alistRes);
        return;
    }
    if (alistRes.startsWith('error404')) {
        const filePath = alistFilePath.substring(alistFilePath.indexOf('/', 1));
        const foldersRes = await fetchAlistPathApi(alistApiPath, '/', alistPwd);
        if (foldersRes.startsWith('error')) {
            r.error(foldersRes);
            r.return(500, foldersRes);
            return;
        }
        const folders = foldersRes.split(',').sort();
        for (let i = 0; i < folders.length; i++) {
            r.warn(`try to fetch alist path from /${folders[i]}${filePath}`);
            const driverRes = await fetchAlistPathApi(alistApiPath, `/${folders[i]}${filePath}`, alistPwd);
            if (!driverRes.startsWith('error')) {
                r.warn(`redirect to: ${driverRes}`);
                r.return(302, driverRes);
                return;
            }
        }
        r.warn(`not found direct ${alistRes}`);
        r.internalRedirect("@backend");
        return;
    }
    r.warn(`not found direct ${alistRes}`);
    r.internalRedirect("@backend");
    return;
}

async function fetchAlistPathApi(alistApiPath, alistFilePath, alistPwd) {
    const alistRequestBody = {
        "path": alistFilePath,
        "password": alistPwd
    }
    try {
        const response = await ngx.fetch(alistApiPath, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json;charset=utf-8'
            },
            max_response_body_size: 65535,
            body: JSON.stringify(alistRequestBody)
        })
        if (response.ok) {
            const result = await response.json();
            if (result === null || result === undefined) {
                return `error: alist_path_api response is null`;
            }
            if (result.message == 'success') {
                if (result.data.type == 'file') {
                    return result.data.files[0].url;
                }
                if (result.data.type == 'folder') {
                    return result.data.files.map(item => item.name).join(',');
                }
            }
            if (result.code == 401) {
                return `error401: alist_path_api ${result.message}`;
            }
            if (result.message.includes('account')) {
                return `error404: alist_path_api ${result.code} ${result.message}`;
            }
            if (result.message == 'file not found' || result.message == 'path not found') {
                return `error404: alist_path_api ${result.message}`;
            }
            return `error: alist_path_api ${result.code} ${result.message}`;
        } else {
            return `error: alist_path_api ${response.status} ${response.statusText}`;
        }
    } catch (error) {
        return (`error: alist_path_api fetchAlistFiled ${error}`);
    }
}

async function fetchEmbyFilePath(itemInfoUri) {
    try {
        const res = await ngx.fetch(itemInfoUri, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json;charset=utf-8',
                'Content-Length': 0,
            },
            max_response_body_size: 65535,
        });
        if (res.ok) {
            const result = await res.json();
            if (result === null || result === undefined) {
                return `error: emby_api itemInfoUri response is null`;
            }
            return result.MediaSources[0].Path;
        } else {
            return (`error: emby_api ${res.status} ${res.statusText}`);
        }
    } catch (error) {
        return (`error: emby_api fetch mediaItemInfo failed,  ${error}`);
    }
}

export default {redirect2Pan};
