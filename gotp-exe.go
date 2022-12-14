package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app := &cli.App{
		Name:  "teleport",
		Usage: "Teleport around your file system!",
		Authors: []*cli.Author{
			{
				Name:  "Brent Mode",
				Email: "brent.a.mode@gmail.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a teleport point",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(cCtx.Args().Slice())
					name, path := cCtx.Args().First(), cCtx.Args().Get(1)
					add(name, path)
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list all teleport points",
				Action: func(cCtx *cli.Context) error {
					list()
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "remove a teleport point",
				Action: func(cCtx *cli.Context) error {
					name := cCtx.Args().First()
					remove(name)
					return nil
				},
			},
			{
				Name:  "goto",
				Usage: "go to a teleport point",
				Action: func(cCtx *cli.Context) error {
					name := cCtx.Args().First()
					res, err := goTo(name)
					check(err)
					fmt.Println(res)
					os.Exit(2)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const TPFILE string = "/.go-teleport"

type TpPoint struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkFileExists() bool {
	home, _ := os.UserHomeDir()
	_, err := os.Stat(home + TPFILE)
	return err == nil
}

func createTpFile() *os.File {
	log.Println("in createTpFile()")
	home, err := os.UserHomeDir()
	check(err)
	f, err := os.Create(home + "/.go-teleport")
	log.Println("file created?")
	check(err)
	jsonHome, err := json.Marshal([]TpPoint{{Name: "home", Path: home}})
	log.Println("json:", string(jsonHome))
	check(err)
	_, err = f.Write(jsonHome)
	check(err)
	return f
}

func createFileIfNotExists() {
	log.Println("checking if file exists")
	if !checkFileExists() {
		log.Println("file did not exist")
		createTpFile()
	}
}

func list() {
	home, _ := os.UserHomeDir()
	if !checkFileExists() {
		createTpFile()
	}

	data, err := os.ReadFile(home + TPFILE)
	check(err)
	tpPoints := &[]TpPoint{}
	err = json.Unmarshal(data, tpPoints)
	check(err)
	fmt.Println("Name [Path]")
	fmt.Println()

	for _, tpPoint := range *tpPoints {
		fmt.Printf("%s [%s]\n", tpPoint.Name, tpPoint.Path)
	}
}

func add(name string, p string) {
	home, _ := os.UserHomeDir()
	log.Println("home:", home)
	// create file if it does not exist yet
	createFileIfNotExists()

	log.Println("user path:", p)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("the path entered (%s) is not a directory\nyou can only teleport to directories", p)
	}
	log.Println("sanitized path:", p)

	// read the file
	data, err := os.ReadFile(home + TPFILE)
	check(err)
	tpPoints := &[]TpPoint{}
	err = json.Unmarshal(data, tpPoints)
	check(err)
	log.Println("points:", tpPoints)

	// add the newest point to the list
	*tpPoints = append(*tpPoints, TpPoint{Name: name, Path: p})
	log.Println("updated points:", tpPoints)

	// write everything to the file
	writeToFile(home, tpPoints)
}

func remove(name string) {
	home, _ := os.UserHomeDir()
	if !checkFileExists() {
		log.Fatal("teleport is not currently set up with a file")
	}

	// read the file
	data, err := os.ReadFile(home + TPFILE)
	check(err)
	tpPoints := &[]TpPoint{}
	err = json.Unmarshal(data, tpPoints)
	check(err)
	log.Println("points:", tpPoints)

	newTpPoints := []TpPoint{}
	for _, tp := range *tpPoints {
		if tp.Name != name {
			newTpPoints = append(newTpPoints, tp)
		} else {
			log.Println("removed point:", tp)
		}
	}

	writeToFile(home, &newTpPoints)
}

func goTo(name string) (string, error) {
	home, _ := os.UserHomeDir()
	if !checkFileExists() {
		log.Fatal("teleport is not currentl set up with a file")
	}

	// read the file
	data, err := os.ReadFile(home + TPFILE)
	check(err)
	tpPoints := &[]TpPoint{}
	err = json.Unmarshal(data, tpPoints)
	check(err)

	for _, tp := range *tpPoints {
		if tp.Name == name {
			return tp.Path, nil
		}
	}

	log.Fatal("the point ", name, " does not exist!")
	return "", errors.New("teleport point does not exist")
}

func writeToFile(home string, points *[]TpPoint) error {
	f, err := os.Create(home + TPFILE)
	if err != nil {
		return err
	}

	jsonPoints, err := json.Marshal(*points)
	if err != nil {
		return err
	}

	log.Println("to be written:\n", string(jsonPoints))
	_, err = f.Write(jsonPoints)
	if err != nil {
		return err
	}
	return nil
}
