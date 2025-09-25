package main

import (
	"fmt"
	"new-client-notification-bot/config"
)

func main() {
	config.Init()

	cfg := config.NewBotConfig()
	logCfg := config.NewLogConfig()

	fmt.Println(logCfg, cfg)

}
