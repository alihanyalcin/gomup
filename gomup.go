package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	"'{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}} {{.Version}} {{.Update.Version}}{{end}}'",
	"-m",
	"all",
}

type Dependency struct {
	name          string
	version       string
	updateVersion string
}

func findDepencencies(path string) ([]Dependency, error) {
	var dependencyList []Dependency

	cmd := exec.Command("go", args...)
	cmd.Dir = path
	list, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s available minor and patch upgrades. Error: %w", path, err)
	}

	dependencies := strings.Split(string(list), "\n")
	for _, dependency := range dependencies {
		if dependency != "''" && dependency != "" {
			d := strings.Split(dependency, " ")
			dependencyList = append(dependencyList, Dependency{
				name:          d[0],
				version:       d[1],
				updateVersion: d[2],
			})
		}
	}

	return dependencyList, nil
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
				wg.Add(1)
				go func() {
					defer wg.Done()

					modPath := path[:(len(path) - len("/go.mod"))]
					d, err := findDepencencies(modPath)
					if err != nil {
						fmt.Println(err)
						return
					}

					if len(d) > 0 {
						dependencies[modPath] = d
					}
				}()
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()
	s.Stop()

	if len(dependencies) == 0 {
		fmt.Println("everything up-to-date")
	}

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
