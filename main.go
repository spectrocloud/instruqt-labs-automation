package main

import (
	"fmt"
	"os"
)

func main() {
	os.Chdir("/app")

	args := os.Args

	fmt.Println(args)

	if len(args) <= 1 {
		fmt.Println("No valid arguments passed.")
		os.Exit(1)
	}
	switch args[1] {
	case "setup":
		Setup(args[2])
	case "cleanup":
		Cleanup()
	case "webserver":
		Webserver()
	default:
		fmt.Println("No valid arguments passed.")
	}
}
