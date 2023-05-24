
package main

import (
	"flag"
	"fmt"
	"os"
	"net/url"

	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
	installer "github.com/kmcsr/server-installer"
)

var loger logger.Logger

func initLogger(){
	loger = logrus.Logger
	if os.Getenv("DEBUG") == "true" {
		loger.SetLevel(logger.TraceLevel)
	}
	_, err := logger.OutputToFile(loger, "./server-installer.log", os.Stdout)
	if err != nil {
		panic(err)
	}
	installer.SetLogger(loger)
	return
}

var (
	TargetVersion string = "latest"
	ServerType string = ""
	InstallPath string = "."
	ExecutableName string = "minecraft"
)

func parseArgs(){
	flag.StringVar(&TargetVersion, "version", TargetVersion,
		"the version of the server need to be installed, could be [latest latest-snapshot]")
	flag.StringVar(&InstallPath, "output", InstallPath,
		"the path need to be installed")
	flag.StringVar(&ExecutableName, "name", ExecutableName,
		"the executable name, without suffix such as '.sh' or '.jar'")
	flag.Usage = func(){
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s (%s):\n", os.Args[0], installer.PkgVersion)
		fmt.Fprint(out, UsageText)
		fmt.Fprintln(out, "Flags:")
		fmt.Fprintln(out, "  -h, -help")
		fmt.Fprintln(out, "        Show this help page")
		flag.PrintDefaults()
		fmt.Fprintln(out, "Args:")
		fmt.Fprintln(out, "  <server_type> string")
		fmt.Fprintf (out, "        type of the server %v (default %q )\n", installer.GetInstallerNames(), ServerType)
		fmt.Fprintln(out, "  <modpack_file> filepath | URL")
		fmt.Fprintln(out, "        the modpack's local path or an URL. If it's an URL, installer will download the modpack first")
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}
	ServerType = flag.Arg(0)
}

func main(){
	parseArgs()
	initLogger()

	fmt.Println()
	switch ServerType {
	case "modpack":
		if flag.NArg() < 2 {
			flag.Usage()
			loger.Fatal("Missing argument <modpack_file>")
		}
		path := flag.Arg(1)
		if _, err := url.ParseRequestURI(path); err == nil {
			var mpath string
			loger.Infof("Downloading modpack %q ...", path)
			if mpath, err = installer.DefaultHTTPClient.DownloadTmp(path, "server-*.mrpack", 0, nil, -1, nil); err != nil {
				loger.Fatalf("Couldn't download modpack %q: %v", path, err)
			}
			defer os.Remove(mpath)
			path = mpath
		}
		loger.Infof("Loading modpack %q ...", path)
		pack, err := installer.OpenMrpack(path)
		if err != nil {
			loger.Fatalf("Couldn't load modpack %q: %v", path, err)
		}
		err = pack.InstallServer(InstallPath)
		pack.Close()
		if err != nil {
			loger.Fatalf("Install modpack error: %v", err)
		}
		var installed string
		minecraft, mok := pack.Deps["minecraft"]
		if forge, ok := pack.Deps["forge"]; ok {
			installed, err = installer.DefaultForgeInstaller.InstallWithLoader(InstallPath, ExecutableName, minecraft, forge)
		}else if fabric, ok := pack.Deps["fabric-loader"]; ok {
			installed, err = installer.DefaultFabricInstaller.InstallWithLoader(InstallPath, ExecutableName, minecraft, fabric)
		}else if quilt, ok := pack.Deps["quilt-loader"]; ok {
			installed, err = installer.DefaultQuiltInstaller.InstallWithLoader(InstallPath, ExecutableName, minecraft, quilt)
		}else if mok {
			installed, err = installer.VanillaIns.Install(InstallPath, ExecutableName, minecraft)
		}else{
			loger.Warnf("Modpack didn't contain any dependencies")
			fmt.Println("\nServer executable file installed to:")
			fmt.Println("NULL")
			return
		}
		if err != nil {
			loger.Fatalf("Install error: %v", err)
		}
		loger.Infof("installed: %s", installed)
		fmt.Println("\nServer executable file installed to:")
		fmt.Println(installed)
	default:
		loger.Infof("Getting version %q for %s server", TargetVersion, ServerType)
		loger.Infof("Install into %q with name %q", InstallPath, ExecutableName)
		fmt.Println()

		ir, ok := installer.Get(ServerType)
		if !ok {
			loger.Fatalf("Could not found installer for server %q", ServerType)
		}
		installed, err := ir.Install(InstallPath, ExecutableName, TargetVersion)
		if err != nil {
			loger.Fatalf("Install error: %v", err)
		}
		loger.Infof("installed: %s", installed)
		fmt.Println("\nServer executable file installed to:")
		fmt.Println(installed)
	}
}
