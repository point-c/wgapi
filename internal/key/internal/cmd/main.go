package main

import (
	"embed"
	"flag"
	helpers "github.com/point-c/generator-helpers"
	"html/template"
	"io"
	"os"
)

//go:embed *.gotmpl
var templates embed.FS

const mainTemplate = "keys.gotmpl"

var args = struct {
	config  string
	outfile string
	pkg     string
}{
	config:  "keys.yml",
	outfile: helpers.OutputFilename(helpers.EnvGoFile()),
	pkg:     helpers.EnvGoPackage(),
}

func init() {
	flag.StringVar(&args.config, "config", args.config, "events config file")
	flag.StringVar(&args.outfile, "out", args.outfile, "output filename")
	flag.StringVar(&args.pkg, "pkg", args.pkg, "gopackage of output")
	flag.Parse()
}

func main() {
	pr, pw := io.Pipe()
	go execTmpl(pw)
	helpers.Check(os.WriteFile(args.outfile, helpers.Must(helpers.GoFmtReader(pr)), os.ModePerm))
}

func execTmpl(pw io.WriteCloser) {
	defer pw.Close()
	helpers.Check(template.Must(template.New("").ParseFS(templates, "*")).ExecuteTemplate(pw, mainTemplate, struct {
		Package string
		Keys    map[string]string
	}{
		Package: args.pkg,
		Keys:    helpers.Must(helpers.UnmarshalYAML[map[string]string](args.config)),
	}))
}
