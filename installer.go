
package installer

type Installer interface {
	// target == "" means latest
	Install(path, name string, target string)(installed string, err error)
}

var Installers = make(map[string]Installer, 10)

func Get(name string)(installer Installer, ok bool){
	installer, ok = Installers[name]
	return
}

func GetInstallerNames()(installers []string){
	installers = make([]string, 0, len(Installers))
	for name, _ := range Installers {
		installers = append(installers, name)
	}
	return
}

type VersionNotFoundErr struct {
	Version string
}

func (e *VersionNotFoundErr)Error()(string){
	return "Version " + e.Version + " not found"
}

type AssetNotFoundErr struct {
	Version string
	Asset   string
}

func (e *AssetNotFoundErr)Error()(string){
	return "Asset " + e.Asset + " is not found in version " + e.Version
}
