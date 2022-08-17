## 常用shell脚本

- [x] [多选菜单](https://github.com/yqchilde/common-dev/blob/main/shell/multiSelect.sh)
- [x] [单选菜单](https://github.com/yqchilde/common-dev/blob/main/shell/singleSelect.sh)
- [x] [github加速](https://github.com/yqchilde/common-dev/blob/main/shell/githubProxy.sh)

## 案例

### 多选菜单

**多选菜单(无序)：**
```shell
source <(curl -sL https://raw.githubusercontents.com/yqchilde/common-dev/main/shell/multiSelect.sh)
declare -A multi_options=(
    ["简体中文"]=true
    ["繁體中文"]=false
    ["English"]=false
)
multiSelect multi_options
for i in "${!multi_options[@]}"; do
    printf "key: %s, val: %s\n" "$i" "${multi_options[$i]}"
done
```

**多选菜单(有序)：**
```shell
source <(curl -sL https://raw.githubusercontents.com/yqchilde/common-dev/main/shell/multiSelect.sh)
declare -A multi_options=(
    ["简体中文"]=true
    ["繁體中文"]=false
    ["English"]=false
)
order=("简体中文" "繁體中文" "English")
multiSelect multi_options order
for i in "${!multi_options[@]}"; do
    printf "key: %s, val: %s\n" "$i" "${multi_options[$i]}"
done
```

<img src="https://github.com/yqchilde/common-dev/blob/main/shell/tests/multiSelect.gif?raw=true" width="400" height="288" alt="multiSelect"/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://github.com/yqchilde/common-dev/blob/main/shell/tests/multiSelect_order.gif?raw=true" width="400" height="288" alt="multiSelect"/>

### 单选菜单

**单选菜单(无序)：**
```shell
source <(curl -sL https://raw.githubusercontents.com/yqchilde/common-dev/main/shell/singleSelect.sh)
declare -A single_options=(
    ["v1.1.1"]=false
    ["v1.1.2"]=true
    ["v1.1.3"]=false
)
singleSelect single_options
for i in "${!single_options[@]}"; do
    printf "key: %s, val: %s\n" "$i" "${single_options[$i]}"
done
```

**单选菜单(有序)：**
```shell
source <(curl -sL https://raw.githubusercontents.com/yqchilde/common-dev/main/shell/singleSelect.sh)
declare -A single_options=(
    ["v1.1.1"]=false
    ["v1.1.2"]=true
    ["v1.1.3"]=false
)
order=("v1.1.2" "v1.1.3" "v1.1.1")
singleSelect single_options order
for i in "${!single_options[@]}"; do
    printf "key: %s, val: %s\n" "$i" "${single_options[$i]}"
done
```

<img src="https://github.com/yqchilde/common-dev/blob/main/shell/tests/singleSelect.gif?raw=true" width="400" height="288" alt="multiSelect"/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://github.com/yqchilde/common-dev/blob/main/shell/tests/singleSelect_order.gif?raw=true" width="400" height="288" alt="multiSelect"/>

### github加速
脚本依赖于 [FastGithub](https://github.com/dotnetcore/FastGithub)，需要执行以下几个步骤：
1. 将对应系统的二进制文件放在服务器
2. cd /usr/bin && ln -s /fastGithub二进制文件路径 /usr/bin/fastGithub
3. 将以上脚本单独写在文件中或者写入`.bashrc`里