package main

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal"
)

func main() {
	appConfig := config.New()
	internal.Run(appConfig)
}
