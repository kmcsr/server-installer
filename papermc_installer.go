package installer

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type (
	PapermcInstaller struct {
		PaperUrl string
	}

	PapermcVersions struct {
		Pid           string   `json:"project_id"`
		Pname         string   `json:"project_name"`
		VersionsGroup []string `json:"version_groups"`
		PaperVersions []string `json:"versions"`
	}

	PapermcBuilders struct {
		Pid       string `json:"project_id"`
		Pname     string `json:"project_name"`
		TgVersion string `json:"version"`
		Builders  []int  `json:"builds"`
	}
)

var DefaultPapermcInstaller = &PapermcInstaller{
	PaperUrl: "https://api.papermc.io/v2/projects/paper",
}
var _ Installer = DefaultPapermcInstaller

func init() {
	Installers["papermc"] = DefaultPapermcInstaller
}

func (r *PapermcInstaller) Install(path, name string, target string) (installed string, err error) {
	return r.InstallWithLoader(path, name, target, "")
}

func (r *PapermcInstaller) InstallWithLoader(path, name string, target string, loader string) (installed string, err error) {
	if len(loader) == 0 {
		var alreadyFind bool = false
		allVersions, err := r.GetInstallerVersions()
		if err != nil {
			return "", err
		}
		if target == "latest" {
			loader = allVersions[len(allVersions)-1]
			alreadyFind = true
		}
		for i := 0; i < len(allVersions); i += 1 {
			if allVersions[i] == target {
				loader = target
				alreadyFind = true
			}
		}
		if !alreadyFind {
			loger.Info("not find the suitable builder, the version should be included in the following list:")
			for i := 0; i < len(allVersions); i += 1 {
				loger.Info("versions:", allVersions[i])
			}
			return "", &VersionNotFoundErr{target}
		}
	}
	buildNumInt, err := r.GetBuildNumber(loader)
	if err != nil {
		return
	}
	buildNum := strconv.Itoa(buildNumInt)
	ExactDownloadeName := "paper-" + loader + "-" + buildNum + ".jar"
	PapermcInstallerUrl, err := url.JoinPath(r.PaperUrl, "versions", loader, "builds", buildNum, "downloads/"+ExactDownloadeName)
	if err != nil {
		return
	}
	loger.Infof("Getting papermc server installer %s at %q...", ExactDownloadeName, PapermcInstallerUrl)
	var buildJar string
	if buildJar, err = DefaultHTTPClient.DownloadDirect(PapermcInstallerUrl, ExactDownloadeName, downloadingCallback(PapermcInstallerUrl)); err != nil {
		return
	}
	installed, err = r.Runbuilder(buildJar, name, ExactDownloadeName, path)
	if err != nil {
		loger.Info("an error occurred while running the server jar file, but you can still do that manually.")
		loger.Error(err)
	}
	return
}

func (r *PapermcInstaller) ListVersions(snapshot bool) (versions []string, err error) {
	data, err := r.GetInstallerVersions()
	if err != nil {
		return
	}
	for _, v := range data {
		versions = append(versions, v)
	}
	return
}

func (r *PapermcInstaller) GetInstallerVersions() (data []string, err error) {
	link := r.PaperUrl
	var versions PapermcVersions
	err = DefaultHTTPClient.GetJson(link, &versions)
	if err != nil {
		return
	}
	data = versions.PaperVersions
	return data, err
}

func (r *PapermcInstaller) GetBuildNumber(version string) (buildNum int, err error) {
	buildUrl, err := url.JoinPath(r.PaperUrl, "versions", version)
	if err != nil {
		return
	}
	var builders PapermcBuilders
	err = DefaultHTTPClient.GetJson(buildUrl, &builders)
	if err != nil {
		return
	}
	buildNum = builders.Builders[len(builders.Builders)-1]
	return buildNum, err
}

func (r *PapermcInstaller) Runbuilder(buildJar string, name string, ExactDownloadName string, path string) (installed string, err error) {
	const SUFFIX string = ".jar"
	serverDirectory := filepath.Join(".", "server-"+ExactDownloadName[0:len(ExactDownloadName)-len(SUFFIX)])
	os.RemoveAll(serverDirectory)
	err = os.MkdirAll(serverDirectory, os.ModePerm)
	if err != nil {
		return
	}
	err = os.Rename(buildJar, filepath.Join(serverDirectory, name+SUFFIX))
	if err != nil {
		return
	}
	buildJar = filepath.Join(serverDirectory, name+SUFFIX)
	loger.Info("Server jar file is successfully installed in path: " + buildJar)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	javapath, err := lookJavaPath()
	if err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, javapath, "-jar", buildJar)
	cmd.Dir = filepath.Join(path, "server-"+ExactDownloadName[0:len(ExactDownloadName)-len(SUFFIX)])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	loger.Infof("Running %q...", cmd.String())
	if err = cmd.Run(); err != nil {
		return
	}
	installed = buildJar
	return
}
