
- English
- [简体中文](./README_zh.MD)

# Minecraft Server Installer

*If this project is useful, please give a star :)*

[Releases Page](https://github.com/kmcsr/server-installer/releases/)

## Introduction

This is a minecraft server installer

You can use this cli to easily install minecraft server.
We also support fabric,forge,spigot

## Flags

```
Flags:
  -h, -help
        show this help message
  -name string
        the executable name, without suffix such as '.sh' or '.jar' (default "minecraft")
  -output string
        the path need to be installed (default ".")
  -server string
        type of the server [fabric,forge,spigot,vanilla]  (default "vanilla")
  -version string
        the version of the server need to be installed, default is the latest (default "latest")
```

## Examples

```sh
# Install minecraft 1.7.10 vanilla server into minecraft.jar
minecraft_installer -name minecraft -version 1.7.10
```

```sh
# Install minecraft 1.19.2 forge server into current directory and the executable is minecraft_server.sh or minecraft_server.bat for windows
minecraft_installer -name minecraft_server -version 1.19.2 -server forge
# Hint: forge installer will make run scripts for the minecraft version that higher or equal than 1.17
#       for version that less than 1.17, you still need to use 'java -jar' to run the server

# Install minecraft 1.16.5 forge server into minecraft_server.jar
minecraft_installer -name minecraft_server -version 1.16.5 -server forge
```

```sh
# Install minecraft 1.19.2 fabric server into server/minecraft_server.jar
minecraft_installer -name minecraft_server -version 1.19.2 -server fabric -path server
```

## TODO

- [ ] PaperMC
- [ ] List server versions
- [ ] Configurable proxy
