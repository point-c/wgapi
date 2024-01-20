package main

import (
	"embed"
	"flag"
	helpers "github.com/point-c/generator-helpers"
	"html/template"
)

//go:embed *.gotmpl
var templates embed.FS

const (
	mainTemplate = "keys.gotmpl"
	testTemplate = "keys_test.gotmpl"
)

var args = struct {
	config string
	out    string
	tout   string
	pkg    string
}{
	config: "keys.yml",
	out:    helpers.OutputFilename(helpers.EnvGoFile()),
	tout:   helpers.TestOutputFilename(helpers.EnvGoFile()),
	pkg:    helpers.EnvGoPackage(),
}

func init() {
	flag.StringVar(&args.config, "config", args.config, "events config file")
	flag.StringVar(&args.out, "out", args.out, "output filename")
	flag.StringVar(&args.tout, "tout", args.tout, "test output filename")
	flag.StringVar(&args.pkg, "pkg", args.pkg, "gopackage of output")
}

func main() {
	flag.Parse()
	Generate(&Dot{
		Package: args.pkg,
		Keys:    helpers.Must(helpers.UnmarshalYAML[map[string]string](args.config)),
	})
}

func Generate(d *Dot) {
	tmpl := helpers.NewTemplate[*template.Template](templates, nil)
	helpers.Generate[*template.Template, *Dot](tmpl, d, mainTemplate, args.out)
	helpers.Generate[*template.Template, *Dot](tmpl, d, testTemplate, args.tout)
}

type Dot struct {
	Package string
	Keys    map[string]string
}
