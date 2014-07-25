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

var all = append([]string{"interface{}"}, primitives...)

var context = map[string]interface{}{
	"Primitives": primitives,
	"All":        all,
}

var goFileReg = regexp.MustCompile("(^[^.].*\\.go)\\.template$")
var illegalCharactersReg = regexp.MustCompile("[^a-zA-Z0-9_]")

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defaultTemplateDir := filepath.Join(wd, "templates")
	if _, err = os.Stat(defaultTemplateDir); err != nil {
		defaultTemplateDir = ""
		err = nil
	}
	defaultDestinationDir := wd
	templateDir := flag.String("templateDir", defaultTemplateDir, "Where to look for the templates.")
	destinationDir := flag.String("destinationDir", defaultDestinationDir, "Where to put the generated files.")

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