package installer

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	SpigotInstaller struct {
	}
)

var _ Installer = (*SpigotInstaller)(nil)

func init() {
	Installers["spigot"] = &SpigotInstaller{}
}

const SpigotBuildToolsURI = "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"

func (*SpigotInstaller) Install(path, name string, target string) (installed string, err error) {
	if _, err = exec.LookPath("git"); err != nil {
		return
	}
	var javapath string
	if javapath, err = lookJavaPath(); err != nil {
		return
	}

	if path, err = filepath.Abs(path); err != nil {
		return
	}

	foundVersion := target
	if target == "" || target == "latest" || target == "latest-snapshot" {
		if target == "latest-snapshot" {
			loger.Info("Warn: spigot do not support snapshot version")
		}
		var versions VanillaVersions
		loger.Info("Getting minecraft version manifest...")
		if versions, err = VanillaIns.GetVersions(); err != nil {
			return
		}
		target = versions.Latest.Release
		foundVersion += "(" + target + ")"
	}

	buildDir, err := os.MkdirTemp("", "spigot-build-")
	if err != nil {
		return
	}
	defer os.RemoveAll(buildDir)
	loger.Infof("Getting %q...", SpigotBuildToolsURI)
	buildToolJar := filepath.Join(buildDir, "BuildTools.jar")
	if err = DefaultHTTPClient.Download(SpigotBuildToolsURI, buildToolJar, 0644, nil, -1,
		downloadingCallback(SpigotBuildToolsURI)); err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, javapath, "-jar", "BuildTools.jar", "--compile", "spigot", "--rev", target)
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	loger.Infof("Running %q...", cmd.String())
	if err = cmd.Run(); err != nil {
		loger.Infof("Build failed. Build log moved to BuildTools.log (if exists)")
		os.Rename(filepath.Join(buildDir, "BuildTools.log.txt"), "BuildTools.log")
		return
	}
	installed = filepath.Join(path, name+".jar")
	if err = renameIfNotExist(filepath.Join(buildDir, "spigot-"+target+".jar"), installed); err != nil {
		return
	}
	return
}

func (r *SpigotInstaller) ListVersions(snapshot bool) (versions []string, err error) {
	versions = []string{"latest"}
	return
}
