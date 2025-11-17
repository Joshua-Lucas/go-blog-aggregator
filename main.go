package main

import (
	"fmt"
	"github.com/Joshua-Lucas/blog-aggregator/internal/config"
	"os"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = cfg.SetUser("Josh")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg, err = config.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v", cfg)

}
