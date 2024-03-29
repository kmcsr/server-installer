
- [English](./README.MD)
- 简体中文

# Minecraft Server Installer

*如果本项目有用, 请给个star吧 :)*

[发布界面](https://github.com/kmcsr/server-installer/releases/)

## 介绍

这是一个需要联网的minecraft服务端下载器

你可以使用本程序简单的安装minecraft及其附属服务端.  
我们也支持一些修改服务端以及一些格式的整合包

| 服务器类型    | 支持  |
|-------------|-------|
| 原版         | 是    |
| Fabric      | 是    |
| Forge       | 是    |
| Quilt       | 是    |
| Spigot      | 是    |
| PaperMC     | 是    |
| ArcLight    | 是    |
| Mohist      | 计划中/否    |
| Catserver   | 计划中/否    |

| 整合包类型     | 支持     |
|--------------|----------|
| Modrinth     | 是       |
| Curseforge   | 计划中/否 |

> 警告: 对于spigot服务端, 您**必须预先安装合适的openjdk**(不仅仅是jre)以及git.  
>       见<https://www.spigotmc.org/wiki/buildtools/#prerequisites>

## 选项

```
Flags:
  -h, -help
        显示这条描述信息
  -name string
        可执行文件名称, 不包含可能的后缀例如'.sh'或'.jar' (默认 "minecraft")
  -output string
        服务端目标安装位置 (默认 ".")
  -version string
        将要安装的minecraft版本, latest或留空为可用的最新版 (默认 "latest")
Args:
  <server_type> string
        服务端类型 [fabric forge spigot vanilla]  (默认 "vanilla")
  <modpack_file> filepath | URL
        整合包的本地路径或URL. 如果为URL则会先将整合包下载到临时路径
```

## 使用示例

### 安装纯服务端

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
# 将 minecraft 1.19.2 fabric服务端下载到 server/minecraft_server.jar
minecraft_installer -name minecraft_server -version 1.19.2 -server fabric -path server
```

```sh
# 将 papermc 1.14.4 服务端下载到 server/minecraft_server.jar
minecraft_installer -name minecraft_server -version 1.14.4 -server papermc
# 注：papermc以及arclight服务器因build区别会区分到不同文件夹中，执行中会自动建立并命名该server文件夹，更加方便识别
```

### 安装整合包

```sh
# 从本地文件安装整合包
minecraft_installer -name modpack_server modpack /path/to/modrinth-modpack.mrpack
# 提示: 目前仅支持modrinth的整合包
```

```sh
# 从网络下载整合包并安装
minecraft_installer -name modpack_server modpack 'https://cdn-raw.modrinth.com/data/sl6XzkCP/versions/i4agaPF2/Automation%20v3.3.mrpack'
# 提示: 如果想要从网络安装,
#       则必须添加前缀 [https://, http://]
```

### 列出服务端可用版本

```sh
minecraft_installer versions
```

```sh
minecraft_installer -version snapshot versions
```
## 计划中

- [ ] 从modrinth搜寻整合包
- [ ] 代理相关
