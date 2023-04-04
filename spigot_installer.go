
package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	SpigotInstaller struct {
	}
)

var _ Installer = (*SpigotInstaller)(nil)

func init(){
	Installers["spigot"] = &SpigotInstaller{}
}

const SpigotBuildToolsURI = "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"

func (*SpigotInstaller)Install(path, name string, target string)(installed string, err error){
	if path, err = filepath.Abs(path); err != nil {
		return
	}

	foundVersion := target
	if target == "" || target == "latest" || target == "latest-snapshot" {
		var versions VanillaVersions
		fmt.Println("Getting minecraft version manifest...")
		if versions, err = VanillaIns.GetVersions(); err != nil {
			return
		}
		if target == "latest-snapshot" {
			fmt.Println("Warn: spigot do not support snapshot version")
		}
		target = versions.Latest.Release
		foundVersion += "(" + target + ")"
	}

	buildDir, err := os.MkdirTemp("", "spigot-build-")
	if err != nil {
		return
	}
	defer os.RemoveAll(buildDir)
	fmt.Printf("Getting %q...\n", SpigotBuildToolsURI)
	var resp *http.Response
	if resp, err = http.DefaultClient.Get(SpigotBuildToolsURI); err != nil {
		return
	}
	defer resp.Body.Close()
	{
		fmt.Printf("Downloading %q...\n", SpigotBuildToolsURI)
		var fd *os.File
		if fd, err = os.Create(filepath.Join(buildDir, "BuildTools.jar")); err != nil {
			return
		}
		_, err = io.Copy(fd, resp.Body)
		resp.Body.Close()
		fd.Close()
		if err != nil {
			return
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	javapath, err := lookJavaPath()
	if err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, javapath, "-jar", "BuildTools.jar", "--compile", "spigot", "--rev", target)
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	fmt.Printf("Running %q...\n", cmd.String())
	if err = cmd.Run(); err != nil {
		fmt.Printf("Build failed. Build log moved to BuildTools.log (if exists)")
		os.Rename(filepath.Join(buildDir, "BuildTools.log.txt"), "BuildTools.log")
		return
	}
	installed = filepath.Join(path, name + ".jar")
	if err = renameIfNotExist(filepath.Join(buildDir, "spigot-" + target + ".jar"), installed); err != nil {
		return
	}
	return
}
