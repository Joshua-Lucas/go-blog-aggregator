package main

import (
	"fmt"
	"github.com/Joshua-Lucas/blog-aggregator/internal"
	"os"
)

func main() {
	config, err := Config.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = config.SetUser("Josh")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	config, err = Config.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v", config)

}
