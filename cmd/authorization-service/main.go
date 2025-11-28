package main

import (
	"authorization-service/internal/config"
	"fmt"
)

func main() {
	cfg, _ := config.MustLoad()
	fmt.Println(*cfg)
}
