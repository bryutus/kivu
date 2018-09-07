package main

import (
	"os"

	"github.com/bryutus/kivu/commands"
	"github.com/urfave/cli"
)

// AppHelpTemplate is the text template for the custom help topic.
var AppHelpTemplate = `NAME:
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
	{{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
	{{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{end}}{{range .VisibleCommands}}
	{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
	{{.Copyright}}{{end}}
`

func main() {
	newApp().Run(os.Args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.CustomAppHelpTemplate = AppHelpTemplate
	app.Name = "kivu"
	app.Usage = "Select npm run script."
	app.Version = "0.1.0"
	app.Author = "bryutus"
	app.Email = "bryutus@gmail.com"
	app.Action = commands.Action
	return app
}
