## 常用shell脚本

- [x] [多选菜单](https://github.com/yqchilde/common-dev/blob/main/shell/multiSelect.sh)
- [x] [单选菜单](https://github.com/yqchilde/common-dev/blob/main/shell/singleSelect.sh)

## 案例

**多选菜单：**
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

<img src="https://github.com/yqchilde/common-dev/blob/main/shell/tests/multiSelect.gif?raw=true" width="450" height="345" alt="multiSelect"/><br/>

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

