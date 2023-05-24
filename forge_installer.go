
package installer

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	ForgeInstaller struct {
		MavenUrl string // Default is "https://maven.minecraftforge.net"
	}
)

var DefaultForgeInstaller = &ForgeInstaller{
	MavenUrl: "https://maven.minecraftforge.net",
}
var _ Installer = DefaultForgeInstaller

func init(){
	Installers["forge"] = DefaultForgeInstaller
}

var v1_17 = Version{
	Major: 1,
	Minor: 17,
	Patch: 0,
}

func (r *ForgeInstaller)Install(path, name string, target string)(installed string, err error){
	return r.InstallWithLoader(path, name, target, "")
}

func (r *ForgeInstaller)InstallWithLoader(path, name string, target string, loader string)(installed string, err error){
	foundVersion := target
	if target == "" || target == "latest" || target == "latest-snapshot" {
		if target == "latest-snapshot" {
			loger.Warn("forge do not support snapshot version")
		}
		var versions VanillaVersions
		loger.Info("Getting minecraft version manifest...")
		if versions, err = VanillaIns.GetVersions(); err != nil {
			return
		}
		target = versions.Latest.Release
		foundVersion += "(" + target + ")"
	}

	var lessV1_17 bool
	{
		var v Version
		if v, err = VersionFromString(target); err != nil {
			return
		}
		lessV1_17 = v.Less(v1_17)
	}

	var version string
	if len(loader) == 0 {
		version, err = r.GetLatestInstaller(target)
		if err != nil {
			return
		}
	}else{
		version = target + "-" + loader
	}
	forgeInstallerUrl, err := url.JoinPath(r.MavenUrl, "net/minecraftforge/forge", version, "forge-" + version + "-installer.jar")
	if err != nil {
		return
	}
	loger.Infof("Getting forge server installer %s at %q...", foundVersion, forgeInstallerUrl)
	var installerJar string
	if installerJar, err = DefaultHTTPClient.DownloadTmp(forgeInstallerUrl, "forge-installer-*.jar", 0644, nil, -1,
		downloadingCallback(forgeInstallerUrl)); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	javapath, err := lookJavaPath()
	if err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, javapath, "-jar", installerJar, "--installServer")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	loger.Infof("Running %q...", cmd.String())
	if err = cmd.Run(); err != nil {
		return
	}

	if lessV1_17 { // < 1.17 use forge-<minecraft_version>-<loader_version>.jar
		installed = filepath.Join(path, name + ".jar")
		if err = renameIfNotExist("forge-" + version + ".jar", installed); err != nil {
			return
		}
		return
	}
	// >= 1.17 use run.sh or run.bat
	installedSh := filepath.Join(path, name + ".sh")
	if err = renameIfNotExist("run.sh", installedSh); err != nil {
		return
	}
	installedBat := filepath.Join(path, name + ".bat")
	if err = renameIfNotExist("run.bat", installedBat); err != nil {
		return
	}
	installed = installedSh
	if runtime.GOOS == "windows" {
		installed = installedBat
	}
	return
}

func (r *ForgeInstaller)GetInstallerVersions()(data MavenMetadata, err error){
	link, err := url.JoinPath(r.MavenUrl, "net/minecraftforge/forge")
	if err != nil {
		return
	}
	return GetMavenMetadata(link)
}

func (r *ForgeInstaller)GetLatestInstaller(target string)(version string, err error){
	data, err := r.GetInstallerVersions()
	if err != nil {
		return
	}
	for _, v := range data.Versioning.Versions {
		if strings.HasPrefix(v, target + "-") {
			version = v
			break
		}
	}
	if len(version) == 0 {
		return "", &VersionNotFoundErr{ "forge-" + target }
	}
	return 
}
