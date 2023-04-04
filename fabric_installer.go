
package installer

import (
	"fmt"
	"net/http"
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

var _ Installer = (*FabricInstaller)(nil)

func init(){
	Installers["fabric"] = &FabricInstaller{
		MetaUrl: "https://meta.fabricmc.net",
	}
}

const fabricServerLauncherProfile = "fabric-server-launcher.properties"
const fabricServerLauncherLink = "https://meta.fabricmc.net/v2/versions/loader/%s/stable/stable/server/jar"

func (r *FabricInstaller)Install(path, name string, target string)(installed string, err error){
	foundVersion := target
	if target == "" || target == "latest" || target == "latest-snapshot" {
		var versions VanillaVersions
		fmt.Println("Getting minecraft version manifest...")
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

	serverLauncherUrl := fmt.Sprintf(fabricServerLauncherLink, target)
	fmt.Printf("Getting fabric server launcher %s at %q...\n", foundVersion, serverLauncherUrl)
	var resp *http.Response
	if resp, err = http.DefaultClient.Get(serverLauncherUrl); err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Downloading %q...\n", serverLauncherUrl)
	installed = filepath.Join(path, name + ".jar")
	if err = safeDownload(resp.Body, installed); err != nil {
		return
	}
	return installed, nil
}

func (r *FabricInstaller)GetInstallers()(res []FabricInstallerVersion, err error){
	tg, err := url.JoinPath(r.MetaUrl, "v2", "versions", "installer")
	if err != nil {
		return
	}
	if err = getHttpJson(tg, &res); err != nil {
		return
	}
	return
}
