# gpp

基于[sing-box](https://github.com/SagerNet/sing-box)的加速器，使用golang编写，支持windows、linux、macos

- 使用vless+ws协议传输
- 支持tcp、udp
- http分流
- gui客户端
- 基于tun代理

# 下载

# 服务端

[下载服务端](https://github.com/danbai225/gpp/releases)

# 客户端

[下载客户端](https://danbai.lanzouq.com/b0064z1wuf)

密码:gpi8
# 编译

## 编译服务端

使用`golang`编译 `cmd/gpp/main.go`获得服务端可执行文件。

## 编译GUI客户端

gui的客户端需要自建构建，需要安装`wails`、`npm`和`golang`，安装方法如下

- 安装`golang`，[下载地址](https://golang.org/dl/)
- 安装`npm` [下载地址](https://nodejs.org/en/download/)
- 安装`wails`，`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

使用`wails`编译

```
wails build
```

# 快速安装服务端脚本（仅支持linux）

```bash
wget https://raw.githubusercontent.com/danbai225/gpp/main/server/install.sh
chmod +x install.sh
./install.sh
```