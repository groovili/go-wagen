package main

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/gobuffalo/packr/v2"
)

const (
	permMode = 0764
)

type application struct {
	Path          string
	Name          string
	Logger        string
	LoggerPackage string
}

func (a *application) String() string {
	return a.Name
}

type dir struct {
	Name  string
	Files map[string]string // [name]template
}

type structure struct {
	Box         *packr.Box
	App         *application
	Files       map[string]string // [name]template
	Directories []*dir
}

func (s *structure) create() {
	for file, temp := range s.Files {
		s.fileFromTemplate(file, temp, s.App)
	}

	for _, d := range s.Directories {
		s.makeDir(d.Name)

		for file, temp := range d.Files {
			s.fileFromTemplate(file, temp, s.App)
		}
	}
}

func (s *structure) makeDir(path string) {
	err := os.MkdirAll(path, permMode)
	if err != nil {
		printErr(err.Error())
		os.Exit(1)
	}

	printSuccess(fmt.Sprintf("Created %s", path))
}

func (s *structure) fileFromTemplate(filePath, templatePath string, args interface{}) {
	f, err := os.Create(filePath)
	if err != nil {
		printErr(err.Error())
		os.Exit(1)
	}
	defer f.Close()

	txt, err := s.Box.FindString(templatePath)
	if err != nil {
		printErr(err.Error())
		os.Exit(1)
	}

	tpl, err := template.New(path.Base(templatePath)).Parse(txt)
	if err != nil {
		printErr(err.Error())
		os.Exit(1)
	}

	if err = tpl.Execute(f, args); err != nil {
		printErr(err.Error())
		os.Exit(1)
	}

	printSuccess(fmt.Sprintf("Created %s", filePath))
}
