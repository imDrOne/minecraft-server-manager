package main

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal"
)

func main() {
	cfg := config.New()
	internal.MigrateUp(cfg)
}
