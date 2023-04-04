
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	installer "github.com/kmcsr/server-installer"
)

var (
	TargetVersion string = "latest"
	ServerType string = "vanilla"
	InstallPath string = "."
	ExecutableName string = "minecraft"
)

func parseArgs(){
	flag.StringVar(&ServerType, "server", ServerType,
		"type of the server [" + strings.Join(installer.GetInstallerNames(), ",") + "] ")
	flag.StringVar(&TargetVersion, "version", TargetVersion,
		"the version of the server need to be installed, default is the latest")
	flag.StringVar(&InstallPath, "output", InstallPath,
		"the path need to be installed")
	flag.StringVar(&ExecutableName, "name", ExecutableName,
		"the executable name, without suffix such as '.sh' or '.jar'")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		fmt.Fprint(out, UsageText)
		fmt.Fprintln(out, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main(){
	parseArgs()

	fmt.Printf(`
Get version %q for %s server
Install into %q with name %q

`, TargetVersion, ServerType, InstallPath, ExecutableName)

	ir, ok := installer.Get(ServerType)
	if !ok {
		fmt.Printf("Error: Could not found installer for server %q\n", ServerType)
		os.Exit(1)
	}
	installed, err := ir.Install(InstallPath, ExecutableName, TargetVersion)
	if err != nil {
		fmt.Printf("Error: Install error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nServer executable file installed to: %s\n", installed)
}
