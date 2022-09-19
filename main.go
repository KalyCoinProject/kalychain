package main

import (
	_ "embed"

	"github.com/KalyCoinProject/kalychain/command/root"
	"github.com/KalyCoinProject/kalychain/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
