
package main

const UsageText = 
`Example:
  minecraft_installer -name minecraft_server -version 1.7.10
        Install minecraft 1.7.10 vanilla server into minecraft_server.jar
  minecraft_installer -name minecraft_server -version 1.19.2 -server forge
        Install minecraft 1.19.2 forge server into current directory and the executable is minecraft_server.sh
        Hint: forge installer will make run scripts for the minecraft version that higher or equal than 1.17
              for version that less than 1.17, you still need to use 'java -jar' to run the server
  minecraft_installer -name minecraft_server -version 1.19.2 -server fabric -path server
        Install minecraft 1.19.2 fabric server into server/minecraft_server.jar
`
