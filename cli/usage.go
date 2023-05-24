
package main

const UsageText = `
minecraft_installer [...flags] <server_type>
minecraft_installer [...flags] modpack <modpack_file>

Example:
  Install servers:
    minecraft_installer -name minecraft_server -version 1.7.10 vanilla
        Install minecraft 1.7.10 vanilla server into minecraft_server.jar
    minecraft_installer -name minecraft_server -version 1.19.2 forge
        Install minecraft 1.19.2 forge server into current directory and the executable is minecraft_server.sh
        Hint: forge installer will make run scripts for the minecraft version that higher or equal than 1.17
              for version that less than 1.17, you still need to use 'java -jar' to run the server
    minecraft_installer -name minecraft_server -version 1.19.2 -output server fabric
        Install minecraft 1.19.2 fabric server into server/minecraft_server.jar
  Install modpacks:
    minecraft_installer -name modpack_server modpack /path/to/modrinch-modpack.mrpack
        Install the modpack from local to the current directory
        Hint: Only support modrinch modpack for now, curseforge is in progress
    minecraft_installer -name modpack_server modpack 'https://cdn-raw.modrinth.com/data/sl6XzkCP/versions/i4agaPF2/Automation%20v3.3.mrpack'
        Install the modpack from internet to the current directory
        Hint: if you want to install modpack from the internet,
              you must add the prefixs [https://, http://]
`
