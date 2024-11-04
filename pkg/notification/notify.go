package notification

import "github.com/gen2brain/beeep"

type NotifyParams struct {
	Title   string
	Message string
}

func Notify(title string, message string) {
	_ = beeep.Notify(title, message, "")
}
