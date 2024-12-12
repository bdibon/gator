package main

import (
	"fmt"
	"log"

	"github.com/bdibon/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	cfg.SetUser("boris")

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println(cfg)
}
