package main

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/app"
)

func main() {
	config := config.New()
	app.Run(config)
}
