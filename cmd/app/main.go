package main

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal"
	"github.com/imDrOne/minecraft-server-manager/internal/service/containers"
)

func main() {
	var tmp containers.MinecraftContainerService
	tmp.Start(context.Background())
	appConfig := config.New()
	internal.Run(appConfig)
}
