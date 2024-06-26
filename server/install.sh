#!/bin/bash
echo "欢迎使用 gpp 服务端安装脚本"
read -p "输入安装路径 (默认是 /usr/local/gpp): " INSTALL_PATH
# 设置默认安装路径
INSTALL_PATH=${INSTALL_PATH:-"/usr/local/gpp"}
read -p "请输入监听地址（默认0.0.0.0）: " LISTEN_ADDRESS
LISTEN_ADDRESS=${LISTEN_ADDRESS:-"0.0.0.0"}
read -p "请输入监听端口（默认5123）: " LISTEN_PORT
LISTEN_PORT=${LISTEN_PORT:-"5123"}
read -p "请输入当前服务器的入口IP: " NET_IP
NET_IP=${NET_IP:-"127.0.0.1"}
echo "请选择一个选项："
echo "1) shadowsocks"
echo "2) socks"
echo "3) vless"
read -p "输入选项 (1-3): " input
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
    *)
        echo "无效选项: $input"
          exit 0
        ;;
esac
echo "您选择的协议为: $PROTOCOL"
echo "您输入的监听地址为: $LISTEN_ADDRESS"
echo "您输入的监听端口为: $LISTEN_PORT"
echo "安装路径为: $INSTALL_PATH"
echo "当前服务器的入口IP为: $NET_IP"
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

echo "下载服务端 。。。"

latest_release_url=$(curl -s https://api.github.com/repos/danbai225/gpp/releases/latest | grep "browser_download_url.*_linux_amd64.tar.gz" | cut -d : -f 2,3 | tr -d \")

filename=$(basename $latest_release_url)

echo "Downloading file: $filename"

curl -LO $latest_release_url

echo "Download complete"

echo "Extracting files"

tar -xzf $filename gpp

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
    echo "stm started with pid \$!"
    exit 0
  fi
elif [ "\$1" = "stop" ]; then
  if [ -f "\$pid_file" ]; then
    pid=\$(cat "\$pid_file")
    echo "Stopping stm with pid \$pid"
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
#fmt.Sprintf("gpp://%s@%s:%d/%s", config.Protocol, ipStr, config.Port, config.UUID
result="gpp://$PROTOCOL@$NET_IP:$LISTEN_PORT/$UUID"
encoded_result=$(echo -n $result | base64)
echo "导入Token：${encoded_result}"