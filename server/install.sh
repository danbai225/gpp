#!/bin/bash
# 使用 command -v 检查 uuidgen 是否存在
if ! command -v uuidgen &> /dev/null
then
    echo "错误: uuidgen 未安装。请安装后继续。"
    exit 1
fi
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
    NET_ADDR=$(curl ipv4.ip.sb)
    NET_ADDR="$NET_ADDR:$LISTEN_PORT"
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
  echo "Directory $INSTALL_PATH created"
fi
# 切换到安装目录
cd "$INSTALL_PATH" || exit
echo "Changed to directory: $PWD"
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
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "下载服务端 。。。"

# 动态地拼接下载URL
latest_release_url=$(curl -s https://api.github.com/repos/danbai225/gpp/releases/latest | grep "browser_download_url.*_linux_$ARCH.tar.gz" | cut -d : -f 2,3 | tr -d \")

filename=$(basename $latest_release_url)

echo "Downloading file: $filename"

curl -LO $latest_release_url

echo "Download complete"

echo "Extracting files"

tar -xzf $filename gpp-server
mv gpp-server gpp
echo "Extraction complete"

rm $filename

chmod +x gpp

cat << EOF > run.sh
#!/bin/bash
cd ${INSTALL_PATH}
pid_file="${INSTALL_PATH}/gpp.pid"
log_file="${INSTALL_PATH}/gpp.log"

if [ "\$1" = "start" ]; then
  if [ -f "\$pid_file" ]; then
    echo "Error: the process is already running"
    exit 1
  else
    echo "Starting gpp"
    nohup ${INSTALL_PATH}/gpp > "\$log_file" 2>&1 &
    echo \$! > "\$pid_file"
    echo "gpp started with pid \$!"
    exit 0
  fi
elif [ "\$1" = "stop" ]; then
  if [ -f "\$pid_file" ]; then
    pid=\$(cat "\$pid_file")
    echo "Stopping gpp with pid \$pid"
    kill "\$pid"
    rm "\$pid_file"
    exit 0
  else
    echo "Error: the process is not running"
    exit 1
  fi
else
  echo "Usage: ${INSTALL_PATH}/run.sh [start|stop]"
  exit 1
fi
EOF

chmod +x run.sh

echo "安装完成,请执行 ${INSTALL_PATH}/run.sh start 启动服务端,执行 ${INSTALL_PATH}/run.sh stop 停止服务端"
read -p "请为您的节点取一个名字: " Name
Name=${Name:-"$NET_ADDR"}
echo "入口地址是: $NET_ADDR"
result="gpp://$PROTOCOL@$NET_ADDR/$UUID"

encoded_result=$(echo -n $result | base64)
echo "导入链接：${encoded_result}#$Name"