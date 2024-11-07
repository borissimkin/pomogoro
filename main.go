package main

import (
	"embed"
	"fmt"
	"github.com/borissimkin/pomogoro/pkg/app"
	"github.com/borissimkin/pomogoro/pkg/notification"
	"github.com/borissimkin/pomogoro/pkg/pomodoro"
	"github.com/borissimkin/pomogoro/pkg/router"
	"github.com/borissimkin/pomogoro/pkg/settings"
	tea "github.com/charmbracelet/bubbletea"
	"os"
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
