package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v2"
)

var args = []string{
	"list",
	"-u",
	"-mod=mod",
	"-f",
	"'{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}'",
	"-m",
	"all",
}

type Dependency struct {
	name          string
	version       string
	updateVersion string
}

func findDepencencies(path string, dependencies map[string][]Dependency, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var dependency []Dependency

	cmd := exec.Command("go", args...)
	cmd.Dir = path
	list, err := cmd.Output()
	if err != nil {
		// TODO: add log
	}

	depList := strings.Split(string(list), "\n")
	r := regexp.MustCompile(`'(.+): (.+) -> (.+)'`)
	for _, d := range depList {
		if d != "''" && d != "" {
			d := r.FindStringSubmatch(d)
			dependency = append(dependency, Dependency{
				name:          d[1],
				version:       d[2],
				updateVersion: d[3],
			})
		}
	}

	dependencies[path] = dependency
}

func find(path string) {
	s := spinner.New(spinner.CharSets[36], 250*time.Millisecond)
	s.Prefix = "Starting gomup "
	s.Start()

	var wg sync.WaitGroup
	dependencies := make(map[string][]Dependency)

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Name() == "go.mod" {
				go findDepencencies(path[:(len(path)-7)], dependencies, &wg)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()
	s.Stop()

	for k, v := range dependencies {
		fmt.Println(k, v)
	}
}

func main() {
	var (
		path string
	)

	app := &cli.App{
		Name:  "gomup",
		Usage: "go module dependency upgrader",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage:       "directory path of go projects",
				Required:    true,
				Destination: &path,
			},
		},
		Action: func(c *cli.Context) error {
			find(path)

			return nil
		},
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
