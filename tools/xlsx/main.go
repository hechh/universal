package main

import "universal/tools/xlsx/internal/parse"

func main() {
	if err := parse.OpenFiles("./define.xlsx"); err != nil {
		panic(err)
	}
	if err := parse.Parse(); err != nil {
		panic(err)
	}
}
