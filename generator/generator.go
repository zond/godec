package main

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

var primitives = []string{
	"string",
	"bool",
	"byte",
	"rune",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
	"float32",
	"float64",
	"complex64",
	"complex128",
}

var context = map[string]interface{}{
	"Primitives": primitives,
	"All":        append([]string{"interface{}", "reflect.Value"}, primitives...),
}

var goFileReg = regexp.MustCompile("(^[^.].*\\.go)\\.template$")
var illegalCharactersReg = regexp.MustCompile("[^a-zA-Z0-9_]")

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defaultDir := filepath.Join(wd, "templates")
	if _, err = os.Stat(defaultDir); err != nil {
		defaultDir = ""
		err = nil
	}
	templateDir := flag.String("templateDir", defaultDir, "Where to look for the templates.")
	destinationDir := flag.String("destinationDir", wd, "Where to put the generated files.")

	flag.Parse()

	if *templateDir == "" {
		flag.Usage()
		return
	}

	templates := template.Must(template.New(".").Funcs(template.FuncMap{
		"gofilter": func(s string) string {
			return illegalCharactersReg.ReplaceAllString(s, "_")
		},
	}).ParseGlob(filepath.Join(*templateDir, "*.go.template")))

	for _, tmpl := range templates.Templates() {
		var destinationFile *os.File
		if destinationFile, err = os.Create(filepath.Join(*destinationDir, filepath.Base(goFileReg.FindStringSubmatch(tmpl.Name())[1]))); err != nil {
			panic(err)
		}
		if err = func() (err error) {
			defer destinationFile.Close()
			if err = tmpl.Execute(destinationFile, context); err != nil {
				return
			}
			return
		}(); err != nil {
			panic(err)
		}

	}
}
