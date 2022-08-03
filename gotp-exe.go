package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "teleport",
		Usage: "Teleport around your file system!",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a teleport point",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(cCtx.Args().Slice())
					name, path := cCtx.Args().First(), cCtx.Args().Get(1)
					err := add(name, path)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list all teleport points",
				Action: func(cCtx *cli.Context) error {
					err := list()
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type TpPoint struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func add(name string, path string) error {
	tpPoint := &TpPoint{Name: name, Path: path}
	home, _ := os.UserHomeDir()
	if _, err := os.Stat(home + "/.go-teleport"); err != nil {
		// new file
		f, err := createTpFile()
		if err != nil {
			return err
		}
		defer f.Close()

		if err := appendTpPointToFile(tpPoint, f); err != nil {
			return err
		}
	} else {
		f, err := os.OpenFile(home+"/.go-teleport", os.O_RDWR, os.ModeAppend)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := appendTpPointToFile(tpPoint, f); err != nil {
			return err
		}
	}
	return nil
}

func list() error {
	home, _ := os.UserHomeDir()
	if _, err := os.Stat(home + "/.go-teleport"); err != nil {
		_, err := createTpFile()
		if err != nil {
			return err
		}
		log.Printf("File created.\nNo teleport points have been set.\n")
	} else {
		data, err := os.ReadFile(home + "/.go-teleport")
		if err != nil {
			return err
		}
		tpPoint := &TpPoint{}
		json.Unmarshal(data, tpPoint)
		fmt.Println("Name [Path]")
		fmt.Println()
		fmt.Printf("%s [%s]\n", tpPoint.Name, tpPoint.Path)
	}
	return nil
}

func createTpFile() (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	f, err := os.Create(home + "/.go-teleport")
	return f, err
}

func appendTpPointToFile(tp *TpPoint, f *os.File) error {
	jsonTpPoint, err := json.Marshal(tp)
	if err != nil {
		return err
	}
	_, err = f.Write(jsonTpPoint)
	if err != nil {
		return err
	}
	return nil
}
