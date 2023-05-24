
package installer

import (
	"fmt"
	"net/url"
	"path/filepath"
)

type (
	FabricInstallerVersion struct {
		Url     string `json:"url"`
		Maven   string `json:"maven"`
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	FabricInstaller struct {
		MetaUrl string // Default is "https://meta.fabricmc.net"
	}
)

var DefaultFabricInstaller = &FabricInstaller{
	MetaUrl: "https://meta.fabricmc.net",
}
var _ Installer = DefaultFabricInstaller

func init(){
	Installers["fabric"] = DefaultFabricInstaller
}

const fabricServerLauncherProfile = "fabric-server-launcher.properties"
const fabricServerLauncherLink = "https://meta.fabricmc.net/v2/versions/loader/%s/%s/stable/server/jar"

func (r *FabricInstaller)Install(path, name string, target string)(installed string, err error){
	return r.InstallWithLoader(path, name, target, "")
}

func (r *FabricInstaller)InstallWithLoader(path, name string, target string, loader string)(installed string, err error){
	foundVersion := target
	if target == "" || target == "latest" || target == "latest-snapshot" {
		var versions VanillaVersions
		loger.Info("Getting minecraft version manifest...")
		if versions, err = VanillaIns.GetVersions(); err != nil {
			return
		}
		if target == "latest-snapshot" {
			target = versions.Latest.Snapshot
			foundVersion += "(" + target + ")"
		}else{
			target = versions.Latest.Release
			foundVersion += "(" + target + ")"
		}
	}
	if loader == "" {
		loader = "stable"
	}

	serverLauncherUrl := fmt.Sprintf(fabricServerLauncherLink, target, loader)
	loger.Infof("Getting fabric server launcher %s at %q...", foundVersion, serverLauncherUrl)
	installed = filepath.Join(path, name + ".jar")
	if err = DefaultHTTPClient.Download(serverLauncherUrl, installed, 0644, nil, -1,
		downloadingCallback(serverLauncherUrl)); err != nil {
		return
	}
	return installed, nil
}

func (r *FabricInstaller)ListVersions(snapshot bool)(versions []string, err error){
	vs, err := r.GetInstallers()
	if err != nil {
		return
	}
	for _, v := range vs {
		if v.Stable || snapshot {
			versions = append(versions, v.Version)
		}
	}
	return
}

func (r *FabricInstaller)GetInstallers()(res []FabricInstallerVersion, err error){
	tg, err := url.JoinPath(r.MetaUrl, "v2", "versions", "installer")
	if err != nil {
		return
	}
	if err = DefaultHTTPClient.GetJson(tg, &res); err != nil {
		return
	}
	return
}
