package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"pomogoro/pkg/pomodoro"
	"pomogoro/pkg/router"
	"pomogoro/pkg/settings"
)

func main() {
	r := router.NewRouter()

	routes := []router.Route{
		router.NewRoute("settings", settings.NewModel(&r)),
		router.NewRoute("pomodoro", pomodoro.NewModel(&r)),
	}

	r.SetRoutes(routes)

	p := tea.NewProgram(r.CurrentRoute().Value)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
