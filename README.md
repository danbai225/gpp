# gpp

gpp加速器，让你的加速器支持主机、mac、linux

# 编译命令行

使用`golang`编译 `cmd/gpp/main.go`获得命令行二进制文件。

# 编译GUI客户端

使用`fyne.io`编译
- windows
```
fyne package -os windows -icon logo.png
```
- linux
```
fyne package -os linux -icon logo.png
```