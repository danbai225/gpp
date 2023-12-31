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

# 使用方法

## 配置文件

使用json的格式，配置文件名为`config.json`，放在程序同级目录下。或者通过第二个参数指定配置文件路径
客户端配置文件和服务端配置文件格式字段相同

## 字段描述

- port 监听端口/服务器端口
- addr 监听地址/服务器地址
- uuid 用于认证的uuid

## 配置文件示例

```json
{
  "port": 5123,
  "addr": "127.0.0.1",
  "uuid": "xxxx-xxx-xxx-xxx-xxx"
}
```

## 服务端启动！

例如我有某加速器，我在加速器上选择加速`英雄联盟国际服`

然后我知道lol游戏文件夹中有个exe叫`client.exe`

那么我就可以这样启动服务端将我的服务端文件名修改为`client.exe`，记得放`config.json`到你的服务端同级目录下

然后启动加速器加速如果可以选择加速模式可以选择进程模式（不行的话可以尝试路由模式）。

## 客户端启动！

同样的将配置文件修改后放入客户端可执行文件同级目录,启动后会自动读取同级目录下的`config.json`文件。

启动成功后会看到`启动成功`提示。第一次使用会下载数据需要等待一会。

# 效果图

![img.png](https://v2.cm/2023/11/13/6551d73019b36.png)