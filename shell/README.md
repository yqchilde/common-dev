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

![测试图片](https://github.com/yqchilde/common-dev/blob/main/shell/tests/multiSelect.gif?raw=true)