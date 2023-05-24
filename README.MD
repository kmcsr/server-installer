
- English
- [简体中文](./README_zh.MD)

Minecraft Server Installer
[![License](https://img.shields.io/github/license/kmcsr/server-installer)](https://github.com/kmcsr/server-installer/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/kmcsr/server-installer/build-tag.yml)](https://github.com/kmcsr/server-installer/actions)
[![Downloads](https://img.shields.io/github/downloads/kmcsr/server-installer/total)](https://github.com/kmcsr/server-installer/releases)
========

*If this project is useful, please give a star :)*

[Releases Page](https://github.com/kmcsr/server-installer/releases/)

## Introduction

This is a minecraft server online installer

You can use this cli to easily install minecraft server.
We also support fabric,forge,spigot and some type of modpacks

> Warn: For spigot server, you **must install suitable openjdk** (not only jre) and git.  
>       See <https://www.spigotmc.org/wiki/buildtools/#prerequisites>

## Flags

```
Flags:
  -h, -help
        Show this help page
  -name string
        the executable name, without suffix such as '.sh' or '.jar' (default "minecraft")
  -output string
        the path need to be installed (default ".")
  -version string
        the version of the server need to be installed, default is the latest (default "latest")
Args:
  <server_type> string
        type of the server [fabric forge spigot vanilla] (default "vanilla" )
  <modpack_file> filepath | URL
        the modpack's local path or an URL. If it's an URL, installer will download the modpack first
```

## Examples

### Install pure servers

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

### Install modpacks

```sh
# Install the modpack from local to the current directory
minecraft_installer -name modpack_server modpack /path/to/modrinch-modpack.mrpack
# Hint: Only support modrinch modpack for now, curseforge is in progress
```

```sh
# Install the modpack from internet to the current directory
minecraft_installer -name modpack_server modpack 'https://cdn-raw.modrinth.com/data/sl6XzkCP/versions/i4agaPF2/Automation%20v3.3.mrpack'
# Hint: if you want to install modpack from the internet,
#       you must add the prefixs [https://, http://]
```

## TODO

- [ ] PaperMC
- [ ] Search modpacks from modrinth
- [ ] List server versions
- [ ] Configurable proxy
