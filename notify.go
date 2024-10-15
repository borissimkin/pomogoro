package main

import "github.com/gen2brain/beeep"

func notify(title string, message string) {
	_ = beeep.Notify(title, message, "assets/information.png")
}
