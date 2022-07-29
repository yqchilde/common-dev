#!/bin/bash

# 单选菜单
function singleSelect {
    cursor_blink_on() { printf "\033[?25h"; }       # 显示光标
    cursor_blink_off() { printf "\033[?25l"; }      # 隐藏光标
    cursor_to() { printf "\033[%s;${2:-1}H" "$1"; } # 光标移动到指定行位置
    print_inactive() { printf "%b %b " "$1" "$2"; } # 打印未激活的选项

    # 获取光标所在终端行数
    function get_cursor_row() {
        IFS=';' read -sdRr -p $'\E[6n' ROW _
        echo "${ROW#*[}"
    }

    local -n options=$1
    local -A selected

    # 定义打印顺序，可自定义
    if [ -z "$2" ]; then
        local -a print_order
        for key in "${!options[@]}"; do
            print_order+=("$key")
        done
    else
        local -n temp=$2
        for key in "${temp[@]}"; do
            print_order+=("$key")
        done
    fi

    # 标记选项状态
    for ((i = 0; i < ${#print_order[@]}; i++)); do
        selected[$i]=${print_order[$i]}
        printf "\n"
    done

    # 确定当前屏幕位置以覆盖选项
    cursor_now_row=$(get_cursor_row)
    local last_row=${cursor_now_row}
    local start_row=$((last_row - ${#options[@]}))

    # 确保在read -s期间在ctrl+c上回显光标和输入
    trap "cursor_blink_on; stty echo; printf '\n'; exit" 2
    cursor_blink_off

    # 监听按键
    key_input() {
        local key
        IFS= read -rsn1 key 2>/dev/null >&2
        if [[ $key = "" ]]; then echo enter; fi
        if [[ $key = $'\x20' ]]; then echo space; fi
        if [[ $key = $'\x1b' ]]; then
            read -rsn2 key
            if [[ $key = [A ]]; then echo up; fi
            if [[ $key = [B ]]; then echo down; fi
        fi
    }

    # 设置选项选中状态
    toggle_option() {
        local idx=$1
        for key in "${!selected[@]}"; do
            if [[ $key = "$idx" ]]; then
                options["${selected[$key]}"]=true
            else
                options["${selected[$key]}"]=false
            fi
        done
    }

    # 打印勾选状态选项
    print_options() {
        for ((i = 0; i < ${#print_order[@]}; i++)); do
            local prefix="[ ]"
            if [[ ${options[${print_order[$i]}]} = true ]]; then
                prefix="[\e[38;5;46m>\e[0m]"
            fi

            cursor_to $((start_row + "$i"))
            print_inactive "$prefix" "${print_order[$i]}"
        done
    }

    local active=0
    for key in "${!selected[@]}"; do
        if [[ ${options[${selected[$key]}]} = true ]]; then
            active="$key"
            break
        fi
    done
    while true; do
        toggle_option "$active"
        print_options "$active"

        case $(key_input) in
        enter)
            print_options -1
            break
            ;;
        up)
            ((active--))
            if [ "$active" -lt 0 ]; then active=$((${#options[@]} - 1)); fi
            ;;
        down)
            ((active++))
            if [ "$active" -ge ${#options[@]} ]; then active=0; fi
            ;;
        esac
    done

    # 光标位置恢复正常
    cursor_to "$last_row"
    cursor_blink_on
}
