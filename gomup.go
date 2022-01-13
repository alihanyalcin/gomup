package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

func findGoModules(path string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Name() == "go.mod" {
				cPath := path[:(len(path) - 7)]
				fmt.Println(cPath)

				args := []string{
					"list",
					"-u",
					"-mod=mod",
					"-f",
					"'{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}'",
					"-m",
					"all",
				}

				cmd := exec.Command("go", args...)
				cmd.Dir = cPath
				list, err := cmd.Output()
				if err != nil {
					return err
				}

				split := strings.Split(string(list), "\n")
				_ = regexp.MustCompile(`'(.+): (.+) -> (.+)'`)
				for _, x := range split {
					if x != "''" && x != "" {
						fmt.Println(x)
					}
				}

			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var (
		path    string
		upgrade bool
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
			&cli.BoolFlag{
				Name:        "upgrade-all",
				Aliases:     []string{"u"},
				Value:       false,
				Usage:       "upgrade all dependencies",
				Destination: &upgrade,
			},
		},
		Action: func(c *cli.Context) error {
			if upgrade {
				fmt.Println("TODO: upgrade dependencies")
			}

			findGoModules(path)

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
