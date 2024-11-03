package main

import (
	"embed"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"pomogoro/pkg/notification"
	"pomogoro/pkg/pomodoro"
	"pomogoro/pkg/router"
	"pomogoro/pkg/settings"
)

//go:embed assets/ring.mp3
var ringAsset embed.FS

func main() {
	notification.SoundAsset = ringAsset

	r := router.NewRouter()

	routes := []router.Route{
		router.NewRoute("pomodoro", pomodoro.NewModel(&r)),
		router.NewRoute("settings", settings.NewModel(&r)),
	}

	r.SetRoutes(routes)

	p := tea.NewProgram(r.CurrentRoute().Value)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
