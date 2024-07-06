# gpp

基于[sing-box](https://github.com/SagerNet/sing-box)+[wails](https://github.com/wailsapp/wails)的加速器，使用golang编写，支持windows、linux、macos

- http分流
- gui客户端
- 基于tun代理

[qq交流群936204503](http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=syMCYJm6Isz_yAxUfrQetpNGioUdpdjO&authKey=lkUyXpKkdAzUwOZYq0m%2BH5Y%2FvAU3XegyxWTm5fM1%2BxOZDdBHJUF%2BODVeNg9MraDl&noverify=0&group_code=936204503)
[TG交流群](https://t.me/+3cX2FOX_owA1ODM1)

# 截图

|                                                         |                                                       |
|---------------------------------------------------------|-------------------------------------------------------|
| ![界面截图](https://imgc.cc/2024/07/06/66888d266d829.png)   | ![英雄联盟](https://imgc.cc/2024/07/06/66888d3c49609.png) |
| ![战地2042](https://imgc.cc/2024/07/06/66888d4ea1807.png) | ![绝地求生](https://imgc.cc/2024/07/06/66888d51e610d.png) |


# 使用教程

## 服务的搭建

在优质线路服务器上运行安装脚本
快速安装服务端脚本（仅支持linux_amd64）
```bash
wget https://raw.githubusercontent.com/danbai225/gpp/main/server/install.sh
chmod +x install.sh
./install.sh
```
根据提示安装完成后会输出导入链接

# 运行客户端

[从releases下载](https://github.com/danbai225/gpp/releases)下载对应系统的客户端以管理员身份运行

点击页面上的`Game`或`Http`字样弹出节点列表窗口，在下方粘贴服务端的链接完成节点导入。
在节点列表选择你的加速节点，如何开始加速。

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