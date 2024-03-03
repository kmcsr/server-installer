package installer

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type (
	QuiltInstaller struct {
		MavenUrl string
	}
)

var DefaultQuiltInstaller = &QuiltInstaller{
	MavenUrl: "https://maven.quiltmc.org/repository/release",
}
var _ Installer = DefaultQuiltInstaller

func init() {
	Installers["quilt"] = DefaultQuiltInstaller
}

func (r *QuiltInstaller) Install(path, name string, target string) (installed string, err error) {
	return r.InstallWithLoader(path, name, target, "")
}

func (r *QuiltInstaller) InstallWithLoader(path, name string, target string, loader string) (installed string, err error) {
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
		} else {
			target = versions.Latest.Release
			foundVersion += "(" + target + ")"
		}
	}

	if len(loader) == 0 {
		loader, err = r.GetLatestInstaller()
		if err != nil {
			return
		}
	}
	quiltInstallerUrl, err := url.JoinPath(r.MavenUrl, "org/quiltmc/quilt-installer", loader, "quilt-installer-"+loader+".jar")
	if err != nil {
		return
	}
	loger.Infof("Getting quilt server installer %s at %q...", foundVersion, quiltInstallerUrl)
	var installerJar string
	if installerJar, err = DefaultHTTPClient.DownloadTmp(quiltInstallerUrl, "quilt-installer-*.jar", 0644, nil, -1,
		downloadingCallback(quiltInstallerUrl)); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	javapath, err := lookJavaPath()
	if err != nil {
		return
	}
	cmd := exec.CommandContext(ctx, javapath, "-jar", installerJar, "install", "server", target, "--download-server", "--install-dir="+path)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	loger.Infof("Running %q...", cmd.String())
	if err = cmd.Run(); err != nil {
		return
	}

	// --download-server flag will install vanilla server to server.jar, we need rename it
	if name == "server" { // name collision
		if err = renameIfNotExist("server.jar", "vanilla_server.jar", 0644); err != nil {
			return
		}
		var fd *os.File
		if fd, err = os.Create("quilt-server-launcher.properties"); err != nil {
			return
		}
		if _, err = fd.Write(([]byte)(time.Now().Format("#" + time.UnixDate + "\n"))); err != nil {
			return
		}
		if _, err = fd.Write(([]byte)(`serverJar=vanilla_server.jar`)); err != nil {
			return
		}
	}
	// Quilt use quilt-server-launch.jar, for some reason, the --create-scripts flag won't work
	installed = filepath.Join(path, name+".jar")
	if err = renameIfNotExist("quilt-server-launch.jar", installed, 0644); err != nil {
		return
	}
	return
}

func (r *QuiltInstaller) ListVersions(snapshot bool) (versions []string, err error) {
	data, err := r.GetInstallerVersions()
	if err != nil {
		return
	}
	for _, v := range data.Versioning.Versions {
		versions = append(versions, v)
	}
	return
}

func (r *QuiltInstaller) GetInstallerVersions() (data MavenMetadata, err error) {
	link, err := url.JoinPath(r.MavenUrl, "org/quiltmc/quilt-installer")
	if err != nil {
		return
	}
	return GetMavenMetadata(link)
}

func (r *QuiltInstaller) GetLatestInstaller() (version string, err error) {
	data, err := r.GetInstallerVersions()
	if err != nil {
		return
	}
	if len(data.Versioning.Versions) == 0 {
		return "", &VersionNotFoundErr{"quilt-latest"}
	}
	var v0, v1 Version
	for _, v := range data.Versioning.Versions {
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
