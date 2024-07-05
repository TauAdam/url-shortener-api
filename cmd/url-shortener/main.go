package main

import (
	"fmt"
	"vigilant-octo-spoon/internal/config"
)

func main() {
	cfg := config.MustLoadEnv()
	fmt.Println(cfg)
}
