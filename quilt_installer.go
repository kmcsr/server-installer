
package installer

type (
	QuiltInstaller struct {
		MavenUrl string
	}
)

var DefaultQuiltInstaller = &QuiltInstaller{
	MavenUrl: "https://maven.quiltmc.org",
}
// var _ Installer = DefaultQuiltInstaller

// func init(){
// 	Installers["quilt"] = DefaultQuiltInstaller
// }

// TODO
