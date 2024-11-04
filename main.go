package main

import (
	"embed"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"pomogoro/pkg/app"
	"pomogoro/pkg/notification"
	"pomogoro/pkg/pomodoro"
	"pomogoro/pkg/router"
	"pomogoro/pkg/settings"
)

//go:embed assets
var assets embed.FS

func main() {
	notification.Assets = assets

	r := router.NewRouter()

	routes := []router.Route{
		router.NewRoute(app.MainPageName, pomodoro.NewModel(&r)),
		router.NewRoute(app.SettingsPageName, settings.NewModel(&r)),
	}

	r.SetRoutes(routes)

	p := tea.NewProgram(r.CurrentRoute().Value)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
