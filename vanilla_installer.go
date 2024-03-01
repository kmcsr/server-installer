package installer

import (
	"path/filepath"
	"time"
)

type (
	DownloadInfo struct {
		Sha1 string `json:"sha1"`
		Size int64  `json:"size"`
		Url  string `json:"url"`
	}
	AssetIndex struct {
		DownloadInfo
		Id        string `json:"id"`
		TotalSize int64  `json:"totalSize"`
	}
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	}
	LibraryDownloadInfo struct {
		DownloadInfo
		Path string `json:"path"`
	}
	LibraryDownloads struct {
		Artifact    *LibraryDownloadInfo           `json:"artifact,omitempty"`
		Classifiers map[string]LibraryDownloadInfo `json:"classifiers,omitempty"`
	}
	LibraryRule map[string]any // TODO
	LibraryInfo struct {
		Name      string            `json:"name"`
		Downloads LibraryDownloads  `json:"downloads"`
		Rules     []LibraryRule     `json:"rules,omitempty"`
		Extract   map[string]any    `json:"extract,omitempty"`
		Natives   map[string]string `json:"natives,omitempty"`
	}

	VanillaVersion struct {
		Id                     string                  `json:"id"`
		AssetIndex             AssetIndex              `json:"assetIndex"`
		Assets                 string                  `json:"assets"`
		ComplianceLevel        int                     `json:"complianceLevel"`
		Downloads              map[string]DownloadInfo `json:"downloads"`
		JavaVersion            JavaVersion             `json:"javaVersion"`
		Libraries              []LibraryInfo           `json:"libraries"`
		Logging                map[string]any          `json:"logging"` // TODO
		MainClass              string                  `json:"mainClass"`
		MinecraftArguments     string                  `json:"minecraftArguments"`
		MinimumLauncherVersion int                     `json:"minimumLauncherVersion"`
		ReleaseTime            time.Time               `json:"releaseTime"`
		Time                   time.Time               `json:"time"`
		Type                   string                  `json:"type"`
	}

	VanillaLatestInfo struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	}
	VanillaVersionInfo struct {
		Id          string    `json:"id"`
		Type        string    `json:"type"`
		Url         string    `json:"url"`
		Time        time.Time `json:"time"`
		ReleaseTime time.Time `json:"releaseTime"`
	}

	VanillaVersions struct {
		Latest   VanillaLatestInfo    `json:"latest"`
		Versions []VanillaVersionInfo `json:"versions"`
	}

	VanillaInstaller struct {
		ManifestUrl string // Default is "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	}
)

var _ Installer = (*VanillaInstaller)(nil)

var VanillaIns = &VanillaInstaller{
	ManifestUrl: "https://launchermeta.mojang.com/mc/game/version_manifest.json",
}

func init() {
	Installers["vanilla"] = VanillaIns
}

func (r *VanillaInstaller) Install(path, name string, target string) (installed string, err error) {
	var res VanillaVersions
	loger.Info("Getting minecraft version manifest...")
	if res, err = r.GetVersions(); err != nil {
		return
	}
	foundVersion := target
	if target == "" || target == "latest" {
		target = res.Latest.Release
		foundVersion += "(" + target + ")"
	} else if target == "latest-snapshot" {
		target = res.Latest.Snapshot
		foundVersion += "(" + target + ")"
	}
	for _, v := range res.Versions {
		if v.Id == target {
			var version VanillaVersion
			loger.Infof("Getting minecraft version %q...", v.Url)
			if version, err = r.GetVersion(v.Url); err != nil {
				return
			}
			info, ok := version.Downloads["server"]
			if !ok {
				return "", &AssetNotFoundErr{foundVersion, "server.jar"}
			}
			installed = filepath.Join(path, name+".jar")
			if err = DefaultHTTPClient.Download(info.Url, installed, 0644, nil, -1,
				downloadingCallback(info.Url)); err != nil {
				return
			}
			return installed, nil
		}
	}
	return "", &VersionNotFoundErr{foundVersion}
}

func (r *VanillaInstaller) ListVersions(snapshot bool) (versions []string, err error) {
	vs, err := r.GetVersions()
	if err != nil {
		return
	}
	for _, v := range vs.Versions {
		if v.Type == "release" || snapshot {
			versions = append(versions, v.Id)
		}
	}
	return
}

func (r *VanillaInstaller) GetVersions() (res VanillaVersions, err error) {
	if err = DefaultHTTPClient.GetJson(r.ManifestUrl, &res); err != nil {
		return
	}
	return
}

func (r *VanillaInstaller) GetVersion(url string) (res VanillaVersion, err error) {
	if err = DefaultHTTPClient.GetJson(url, &res); err != nil {
		return
	}
	return
}
