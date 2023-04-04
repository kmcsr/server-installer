
- [English](./README.MD)
- 简体中文

# Minecraft Server Installer

*如果本项目有用, 请给个star吧 :)*

[发布界面](https://github.com/kmcsr/server-installer/releases/)

## 介绍

这是一个minecraft服务端下载器

你可以使用本程序简单的安装minecraft及其附属服务端.
我们目前支持 vanilla,fabric,forge,spigot

## 选项

```
Flags:
  -h, -help
        显示这条描述信息
  -name string
        可执行文件名称, 不包含可能的后缀例如'.sh'或'.jar' (默认 "minecraft")
  -output string
        服务端目标安装位置 (默认 ".")
  -server string
        服务端类型 [fabric,forge,spigot,vanilla]  (默认 "vanilla")
  -version string
        将要安装的minecraft版本, latest或留空为可用的最新版 (默认 "latest")
```

## 使用示例

```sh
# 将原版 minecraft 1.7.10 服务端下载到 minecraft.jar
minecraft_installer -name minecraft -version 1.7.10
```

```sh
# 将 minecraft 1.19.2 forge服务端下载到当前路径下执行脚本将重命名为 minecraft_server.sh 及 windows下的minecraft_server.bat
minecraft_installer -name minecraft_server -version 1.19.2 -server forge
# 提示: forge下载器会为大于等于1.17的minecraft版本创建一个执行脚本, 您应该直接执行该脚本以启动服务端
#       对于小于1.17的版本, 您仍然需要使用 'java -jar' 启动服务端

# 将 minecraft 1.16.5 forge 服务端下载到 minecraft_server.jar
minecraft_installer -name minecraft_server -version 1.16.5 -server forge
```

```sh
# 将 minecraft 1.7.10 fabric服务端下载到 server/minecraft_server.jar
minecraft_installer -name minecraft_server -version 1.19.2 -server fabric -path server
```
