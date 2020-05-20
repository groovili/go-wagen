package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	packr "github.com/gobuffalo/packr/v2"
)

const permMode = 0764

type application struct {
	Path string
	Name string
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

	printMsg(fmt.Sprintf("Created %s", path))
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

	printMsg(fmt.Sprintf("Created %s", filePath))
}

type actionFunc func(answer string) error

type action struct {
	Question string
	Validate actionFunc
	Action   actionFunc
}

func printErr(msg string) {
	fmt.Fprint(os.Stderr, fmt.Sprintf("%s\r\n", msg))
}

func printMsg(msg string) {
	fmt.Fprint(os.Stdin, fmt.Sprintf("%s\r\n", msg))
}

func userAction(a *action) error {
	printMsg(a.Question)

	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		inp := scn.Text()
		if err := a.Validate(inp); err != nil {
			printErr(err.Error())
			continue
		}

		if err := a.Action(inp); err != nil {
			return err
		}

		break
	}

	return nil
}

func main() {
	var pathToApp string
	flag.StringVar(&pathToApp, "path", "", "Specify absolute path to app")
	flag.Parse()

	if len(pathToApp) == 0 {
		printErr("--path is required")
		os.Exit(1)
	}

	if !path.IsAbs(pathToApp) {
		printErr("--path should be absolute")
		os.Exit(1)
	}

	pathToApp = path.Clean(pathToApp)
	sep := string(os.PathSeparator)

	stat, err := os.Stat(pathToApp)
	if err != nil && !os.IsNotExist(err) {
		printErr(err.Error())
		os.Exit(1)
	} else if os.IsNotExist(err) {
		err := userAction(&action{
			Question: fmt.Sprintf("Dir %s doesn't exist. Create? [y/N]:", pathToApp),
			Validate: func(answer string) error {
				a := strings.ToLower(answer)
				if a != "y" && a != "n" {
					return errors.New("Invalid option. Only [y/N] available:")
				}

				return nil
			},
			Action: func(answer string) error {
				if strings.ToLower(answer) == "n" {
					return errors.New("Can't continue without app dir")
				}

				return os.Mkdir(pathToApp, permMode)
			},
		})

		if err != nil {
			printErr(err.Error())
			os.Exit(0)
		}

		printMsg(fmt.Sprintf("Created %s", pathToApp))

		stat, err = os.Stat(pathToApp)
		if err != nil {
			printErr(err.Error())
			os.Exit(1)
		}
	}

	if !stat.IsDir() {
		printErr("--path should be a directory")
		os.Exit(1)
	}

	app := new(application)
	app.Path = fmt.Sprintf("%s%s", pathToApp, sep)

	err = userAction(&action{
		Question: "Enter app name [a-z0-9_]:",
		Validate: func(answer string) error {
			err := errors.New("App name should be in lower snake case [a-z0-9_]")
			if len(answer) == 0 {
				return err
			}
			r := regexp.MustCompile("^[a-z0-9_]*$")
			if !r.MatchString(answer) {
				return err
			}

			return nil
		},
		Action: func(answer string) error {
			app.Name = answer

			return nil
		},
	})
	if err != nil {
		printErr(err.Error())
		os.Exit(1)
	}

	printMsg(fmt.Sprintf("Creating %s application..", app))

	s := structure{
		Box: packr.New("tplBox", "./templates"),
		App: app,
		Files: map[string]string{
			fmt.Sprintf("%s%s", app.Path, "go.mod"):   "go.mod.tmpl",
			fmt.Sprintf("%s%s", app.Path, "Makefile"): "Makefile.tmpl",
		},
		Directories: make([]*dir, 0),
	}

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s%s%s", app.Path, "cmd", sep, app),
		Files: map[string]string{
			fmt.Sprintf("%s%s%s%s%s%s", app.Path, "cmd", sep, app, sep, fmt.Sprintf("%s.go", app)): "app.go.tmpl",
		},
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s", app.Path, "config"),
		Files: map[string]string{
			fmt.Sprintf("%s%s%s%s", app.Path, "config", sep, "app.local.yml"): "config.yml.tmpl",
			fmt.Sprintf("%s%s%s%s", app.Path, "config", sep, "app.dev.yml"):   "config.yml.tmpl",
			fmt.Sprintf("%s%s%s%s", app.Path, "config", sep, "app.yml"):       "config.yml.tmpl",
		},
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s", app.Path, "deploy"),
		Files: map[string]string{
			fmt.Sprintf("%s%s%s%s", app.Path, "deploy", sep, "Dockerfile"):         "Dockerfile.tmpl",
			fmt.Sprintf("%s%s%s%s", app.Path, "deploy", sep, "docker-compose.yml"): "docker-compose.yml.tmpl",
		},
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s", app.Path, "internal"),
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s", app.Path, "vendor"),
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s", app.Path, "storage"),
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s%s%s%s%s", app.Path, "server", sep, "hadlers", sep, "rest"),
	})

	s.Directories = append(s.Directories, &dir{
		Name: fmt.Sprintf("%s%s%s%s", app.Path, "server", sep, "middleware"),
	})

	s.create()

	printMsg("Success!")
}
