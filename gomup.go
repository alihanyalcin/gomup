package main

import (
	"fmt"
	"gomup/ui"
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
	var dependencies []Dependency

	cmd := exec.Command("go", args...)
	cmd.Dir = path
	list, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s available minor and patch upgrades. error: %w", path, err)
	}

	for _, dependency := range strings.Split(string(list), "\n") {
		if dependency != "''" && dependency != "" {
			d := strings.Split(dependency, " ")
			dependencies = append(dependencies, Dependency{
				name:          d[0],
				version:       d[1],
				updateVersion: d[2],
			})
		}
	}

	return dependencies, nil
}

func find(path string) map[string][]Dependency {
	var dependencies = make(map[string][]Dependency)
	var wg sync.WaitGroup

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Name() == "go.mod" {
				wg.Add(1)
				go func() {
					defer wg.Done()

					modPath := path[:(len(path) - len("go.mod"))]
					d, err := findDepencencies(modPath)
					if err != nil {
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

	return dependencies
}

func run(path string) error {
	err := checkPath(path)
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[36], 250*time.Millisecond)
	s.Prefix = "Starting gomup "
	s.Start()

	dependencies := find(path)

	s.Stop()

	if len(dependencies) == 0 {
		fmt.Println("everything up-to-date")
		return nil
	}

	ui.Start()

	for k, v := range dependencies {
		fmt.Println(k, v)
	}

	return nil
}

func checkPath(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
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
			return run(path)
		},
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
