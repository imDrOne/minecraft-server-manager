package main

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/app"
)

func main() {
	cfg := config.New()
	app.Run(cfg)
}
