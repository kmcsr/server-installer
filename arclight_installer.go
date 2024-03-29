package installer

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type (
	ArclightInstaller struct {}

	ArclightAssets struct {
		AssetsUrl   string `json:"url"`
		AssetsName  string `json:"name"`
		DownloadUrl string `json:"browser_download_url"`
	}

	ArclightRelease struct {
		Assets      []ArclightAssets `json:"assets"`
		IsExpired   bool
		PublishTime string `json:"published_at"`
	}
)

var DefaultArclightInstaller = &ArclightInstaller{}

var _ Installer = DefaultArclightInstaller

func init() {
	Installers["arclight"] = DefaultArclightInstaller
}

func (r *ArclightInstaller) Install(path, name string, target string) (installed string, err error) {
	return r.InstallWithLoader(path, name, target, "")
}

func (r *ArclightInstaller) InstallWithLoader(path, name string, target string, loader string) (installed string, err error) {
	versions, err := r.GetInstallerVersions()
	if err != nil {
		return "", err
	}
	if len(loader) == 0 {
		var alreadyFind bool = false
		allVersions := r.GetOnlyVersions(versions)
		if target == "latest" {
			loader, err = r.GetLatestVersion()
			if err != nil {
				return "", err
			}
			alreadyFind = true
		}
		for _, version := range allVersions {
			if version == target {
				loader = target
				alreadyFind = true
			}
		}
		if !alreadyFind {
			loger.Info("not find the suitable builder, the version should be included in the following list:")
			for _, version := range allVersions {
				if versions[version].IsExpired {
					loger.Info("versions:", version, "  EXPIRED, DO NOT SUPPORT")
				} else {
					loger.Info("versions:", version)
				}
			}
			return "", &VersionNotFoundErr{target}
		}
	}
	exactDownloadeName := versions[loader].Assets[0].AssetsName
	ArclightInstallerUrl := versions[loader].Assets[0].DownloadUrl
	if version, ok := versions[loader]; ok && version.IsExpired {
		loger.Fatal("Sorry, the one you choose has already expired, try another version.")
		return "", &VersionNotFoundErr{target}
	}
	var buildJar string
	if buildJar, err = DefaultHTTPClient.DownloadDirect(ArclightInstallerUrl, exactDownloadeName, downloadingCallback(ArclightInstallerUrl)); err != nil {
		return
	}
	installed, err = r.Runbuilder(buildJar, exactDownloadeName, path)
	if err != nil {
		loger.Info("an error occurred while running the server jar file, but you can still do that manually.")
		loger.Error(err)
	}
	return
}

func (r *ArclightInstaller) ListVersions(snapshot bool) ([]string, error) {
	additionalVersions, err := r.GetInstallerVersions()
	if err != nil {
		return nil, err
	}
	return r.GetOnlyVersions(additionalVersions), err
}

func (r *ArclightInstaller) GetLatestVersion() (version string, err error) {
	additionalVersions, err := r.GetInstallerVersions()
	if err != nil {
		return
	}
	var dataVersions []string = r.GetOnlyVersions(additionalVersions)
	var v0, v1 Version
	for _, v := range dataVersions {
		if v1, err = VersionFromString(v); err != nil {
			return
		}
		if v0.Less(v1) {
			v0 = v1
		}
	}
	version = v0.String()
	return
}

func (r *ArclightInstaller) GetInstallerVersions() (map[string]ArclightRelease, error) {
	additionalVersions := make(map[string]ArclightRelease)
	link := "https://api.github.com/repos/IzzelAliz/Arclight/releases"
	var releases []*ArclightRelease
	err := DefaultHTTPClient.GetJson(link, &releases)
	if err != nil {
		return additionalVersions, err
	}
	for _, release := range releases {
		details := strings.Split(release.Assets[0].AssetsName, "-")
		// details should be ["arclight","forge","{VERSION}","{BUILDNUM}.jar"], so append value of index 2
		timeDetails := strings.Split(release.PublishTime, "-")
		// time should be "{YEAR}-{MONTH}-{DATE}T{CLOCK}}"
		year, err := strconv.Atoi(timeDetails[0])
		if err != nil {
			return additionalVersions, err
		}
		month, err := strconv.Atoi(timeDetails[1])
		if err != nil {
			return additionalVersions, err
		}
		if year < 2024 || (year == 2024 && month < 2) {
			release.IsExpired = true
		} else {
			release.IsExpired = false
		}
		if len(additionalVersions[details[2]].Assets) == 0 {
			additionalVersions[details[2]] = *release
		}
		// to get the newest builder for each version
	}
	return additionalVersions, err
}

func (r *ArclightInstaller) GetOnlyVersions(additionalVersions map[string]ArclightRelease) (versions []string) {
	for k, _ := range additionalVersions {
		versions = append(versions, k)
	}
	return
}

func (r *ArclightInstaller) Runbuilder(buildJar string, ExactDownloadName string, path string) (installed string, err error) {
	const SUFFIX_LENGTH int = 4
	NameWithoutSuffix := ExactDownloadName[0 : len(ExactDownloadName)-SUFFIX_LENGTH]
	serverDirectory := filepath.Join(".", "server-"+NameWithoutSuffix)
	os.RemoveAll(serverDirectory)
	err = os.MkdirAll(serverDirectory, os.ModePerm)
	if err != nil {
		return
	}
	err = os.Rename(buildJar, filepath.Join(serverDirectory, ExactDownloadName))
	if err != nil {
		return
	}
	buildJar = filepath.Join(serverDirectory, ExactDownloadName)
	loger.Info("Server jar file is successfully installed in path: " + buildJar)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	javapath, err := lookJavaPath()
	if err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, javapath, "-jar", buildJar)
	cmd.Dir = filepath.Join(path, "server-"+ExactDownloadName[0 : len(ExactDownloadName)-SUFFIX_LENGTH])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	loger.Infof("Running %q...", cmd.String())
	if err = cmd.Run(); err != nil {
		return
	}
	installed = buildJar
	return
}
