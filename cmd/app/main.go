package main

import (
	"github.com/serhiimazurok/ecoflow-status-bot/internal/app"
)

const configDir = "configs"

func main() {
	app.Run(configDir)
}
