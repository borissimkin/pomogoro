package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"pomogoro/pkg/pomodoro"
)

func main() {
	// todo: alt screen?
	p := tea.NewProgram(pomodoro.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
