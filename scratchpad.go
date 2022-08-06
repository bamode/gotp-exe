package main

import (
	"log"
	"os"
)

func amain() {
	fi, err := os.Stat("tester.txt")
	check(err)
	log.Println(fi)
}

type Point struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
