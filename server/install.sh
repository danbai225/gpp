#!/bin/bash
# 检查root权限
if [ "$(id -u)" -ne 0 ]; then
    echo "警告: 非root用户运行，某些功能可能受限。建议使用root或sudo运行此脚本。"
fi

# 检查必要的工具
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "错误: $1 未安装。"
        echo "尝试安装 $1..."
        
        # 检测包管理器并安装
        if command -v apt &> /dev/null; then
            apt update && apt install -y $2
        elif command -v dnf &> /dev/null; then
            dnf install -y $2
        elif command -v yum &> /dev/null; then
            yum install -y $2
        elif command -v zypper &> /dev/null; then
            zypper install -y $2
        elif command -v pacman &> /dev/null; then
            pacman -S --noconfirm $2
        else
            echo "错误: 无法确定系统的包管理器。请手动安装 $1 后继续。"
            exit 1
        fi
        
        # 再次检查安装是否成功
        if ! command -v $1 &> /dev/null; then
            echo "错误: 安装 $1 失败，请手动安装后继续。"
            exit 1
        fi
    fi
}

# 检查必要的命令
check_command "uuidgen" "uuid-runtime"
check_command "curl" "curl"
check_command "tar" "tar"

echo "欢迎使用 gpp 服务端安装脚本"
read -p "输入安装路径 (默认是 /usr/local/gpp): " INSTALL_PATH
# 设置默认安装路径
INSTALL_PATH=${INSTALL_PATH:-"/usr/local/gpp"}
read -p "请输入监听地址（默认0.0.0.0）: " LISTEN_ADDRESS
LISTEN_ADDRESS=${LISTEN_ADDRESS:-"0.0.0.0"}
read -p "请输入监听端口（默认5123）: " LISTEN_PORT
LISTEN_PORT=${LISTEN_PORT:-"5123"}
read -p "请输入你的客户端入口地址(有中转就是中转地址不填默认当前服务器ip+端口): " NET_ADDR
NET_ADDR=${NET_ADDR:-""}
# 如果NET_ADDR变量为空，则获取外网IP地址
if [ -z "$NET_ADDR" ]; then
    NET_ADDR=$(curl -s ipv4.ip.sb || curl -s ifconfig.me || curl -s icanhazip.com)
    if [ -z "$NET_ADDR" ]; then
        echo "警告: 无法自动获取外网IP地址，请手动指定。"
        read -p "请输入你的客户端入口地址: " NET_ADDR
        if [ -z "$NET_ADDR" ]; then
            echo "错误: 未提供入口地址，安装终止。"
            exit 1
        fi
    else
        NET_ADDR="$NET_ADDR:$LISTEN_PORT"
    fi
fi

echo "请选择一个选项："
echo "1) shadowsocks"
echo "2) socks"
echo "3) vless"
echo "4) hysteria2"
read -p "输入选项 (1-4): " input
PROTOCOL="vless"
case $input in
    1)
        PROTOCOL="shadowsocks"
        ;;
    2)
        PROTOCOL="socks"
        ;;
    3)
        PROTOCOL="vless"
        ;;
    4)
        PROTOCOL="hysteria2"
        ;;
    *)
        echo "无效选项: $input"
          exit 0
        ;;
esac
echo "您选择的协议为: $PROTOCOL"
echo "您输入的监听地址为: $LISTEN_ADDRESS"
echo "您输入的监听端口为: $LISTEN_PORT"
echo "安装路径为: $INSTALL_PATH"
echo "您的入口地址: $NET_ADDR"
# 检查目录是否存在，如果不存在则创建
if [ ! -d "$INSTALL_PATH" ]; then
  mkdir -p "$INSTALL_PATH"
  echo "目录 $INSTALL_PATH 已创建"
fi
# 切换到安装目录
cd "$INSTALL_PATH" || exit
echo "已切换到目录: $PWD"
UUID=$(uuidgen)
cat << EOF > config.json
{
  "protocol": "$PROTOCOL",
  "port": $LISTEN_PORT,
  "addr": "$LISTEN_ADDRESS",
  "uuid":"$UUID"
}
EOF

echo "检测系统架构..."

ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l|armv7)
        ARCH="arm"
        ;;
    i386|i686)
        ARCH="386"
        ;;
    *)
        echo "不支持的架构: $ARCH"
        exit 1
        ;;
esac

echo "下载服务端 。。。"

# 动态地拼接下载URL
latest_release_url=$(curl -s https://api.github.com/repos/danbai225/gpp/releases/latest | grep "browser_download_url.*_linux_$ARCH.tar.gz" | cut -d : -f 2,3 | tr -d \")

# 检查是否成功获取URL
if [ -z "$latest_release_url" ]; then
    echo "错误: 无法获取下载URL，请检查网络连接或手动下载。"
    exit 1
fi

filename=$(basename $latest_release_url)

echo "下载文件: $filename"

curl -LO $latest_release_url

if [ $? -ne 0 ]; then
    echo "错误: 下载失败，请检查网络连接或手动下载。"
    exit 1
fi

echo "下载完成"

echo "解压文件"

tar -xzf $filename gpp-server
if [ $? -ne 0 ]; then
    echo "错误: 解压失败。"
    exit 1
fi

mv gpp-server gpp
echo "解压完成"

rm $filename

chmod +x gpp

# 创建运行脚本
cat << EOF > run.sh
#!/bin/bash
cd ${INSTALL_PATH}
pid_file="${INSTALL_PATH}/gpp.pid"
log_file="${INSTALL_PATH}/gpp.log"

if [ "\$1" = "start" ]; then
  if [ -f "\$pid_file" ]; then
    echo "错误: 进程已经在运行中"
    exit 1
  else
    echo "启动 gpp"
    nohup ${INSTALL_PATH}/gpp > "\$log_file" 2>&1 &
    echo \$! > "\$pid_file"
    echo "gpp 已启动，进程ID为 \$!"
    exit 0
  fi
elif [ "\$1" = "stop" ]; then
  if [ -f "\$pid_file" ]; then
    pid=\$(cat "\$pid_file")
    echo "停止 gpp，进程ID为 \$pid"
    kill "\$pid"
    rm "\$pid_file"
    exit 0
  else
    echo "错误: 进程未运行"
    exit 1
  fi
elif [ "\$1" = "restart" ]; then
  \$0 stop
  sleep 1
  \$0 start
elif [ "\$1" = "status" ]; then
  if [ -f "\$pid_file" ]; then
    pid=\$(cat "\$pid_file")
    if ps -p \$pid > /dev/null; then
      echo "gpp 正在运行，进程ID为 \$pid"
    else
      echo "gpp 似乎已崩溃，进程ID \$pid 不存在"
      rm "\$pid_file"
    fi
  else
    echo "gpp 未运行"
  fi
else
  echo "用法: ${INSTALL_PATH}/run.sh [start|stop|restart|status]"
  exit 1
fi
EOF

chmod +x run.sh

# 检测是否支持systemd
HAS_SYSTEMD=false
if command -v systemctl &> /dev/null && [ ! -f /.dockerenv ]; then
    HAS_SYSTEMD=true
fi

# 如果支持systemd，创建systemd服务文件
if [ "$HAS_SYSTEMD" = true ]; then
    echo "检测到系统支持systemd，创建系统服务..."
    
    cat << EOF > /etc/systemd/system/gpp.service
[Unit]
Description=GPP Proxy Service
After=network.target

[Service]
Type=simple
User=$(whoami)
WorkingDirectory=${INSTALL_PATH}
ExecStart=${INSTALL_PATH}/gpp
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

    # 重新加载systemd配置
    systemctl daemon-reload
    
    echo "systemd服务已创建。您可以使用以下命令管理服务:"
    echo "启动服务: sudo systemctl start gpp 或 sudo service gpp start"
    echo "停止服务: sudo systemctl stop gpp 或 sudo service gpp stop"
    echo "查看状态: sudo systemctl status gpp 或 sudo service gpp status"
    echo "启用开机自启: sudo systemctl enable gpp"
    
    # 询问是否立即启动服务并设置开机自启
    read -p "是否立即启动服务? (y/n): " START_SERVICE
    if [ "$START_SERVICE" = "y" ] || [ "$START_SERVICE" = "Y" ]; then
        systemctl start gpp
        echo "服务已启动"
    fi
    
    read -p "是否设置开机自启? (y/n): " ENABLE_SERVICE
    if [ "$ENABLE_SERVICE" = "y" ] || [ "$ENABLE_SERVICE" = "Y" ]; then
        systemctl enable gpp
        echo "服务已设置为开机自启"
    fi
else
    echo "使用传统方式管理服务"
    echo "安装完成,请执行 ${INSTALL_PATH}/run.sh start 启动服务端,执行 ${INSTALL_PATH}/run.sh stop 停止服务端"
    
    # 询问是否立即启动服务
    read -p "是否立即启动服务? (y/n): " START_SERVICE
    if [ "$START_SERVICE" = "y" ] || [ "$START_SERVICE" = "Y" ]; then
        ${INSTALL_PATH}/run.sh start
        echo "服务已启动"
    fi
fi

read -p "请为您的节点取一个名字: " Name
Name=${Name:-"$NET_ADDR"}
echo "入口地址是: $NET_ADDR"
result="gpp://$PROTOCOL@$NET_ADDR/$UUID"

# 编码链接（兼容不同系统）
if command -v base64 &> /dev/null; then
    encoded_result=$(echo -n $result | base64 | tr -d '\n')
else
    # 如果没有base64命令，生成未编码链接
    encoded_result=$result
    echo "警告: 未找到base64命令，生成未编码链接"
fi

echo "导入链接：${encoded_result}#$Name"
echo "安装完成！"