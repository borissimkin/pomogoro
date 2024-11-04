package settings

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

const (
	toggleItem kindFormItem = "toggle"
	numberItem kindFormItem = "number"
)

var (
	onStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	offStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
)

type kindFormItem string

type formItem struct {
	title  string
	value  int
	kind   kindFormItem
	limits *limits
}

// todo: убрать эти isToggle сделать на интерфейсах?

func (item *formItem) isToggle() bool {
	return item.kind == toggleItem
}

func (item *formItem) isNumber() bool {
	return item.kind == numberItem
}

func (item *formItem) Enter() {
	if !item.isToggle() {
		return
	}

	if item.value >= 1 {
		item.value = 0
	} else {
		item.value = 1
	}
}

func (item *formItem) Increase() {
	if item.isToggle() {
		item.value = 1
		return
	}

	value := item.value + 1

	if item.limits != nil && item.limits.max < value {
		return
	}

	item.value = value
}

func (item *formItem) Decrease() {
	if item.isToggle() {
		item.value = 0
	}

	value := item.value - 1

	if item.limits != nil && item.limits.min > value {
		return
	}

	item.value = value
}

func (item *formItem) View() string {
	if item.isToggle() {
		return item.toggleItemView()
	}

	if item.isNumber() {
		return item.numberItemView()
	}

	return ""
}

func (item *formItem) toggleItemView() string {
	s := ""

	value := offStyle.Render("off")

	if item.value == 1 {
		value = onStyle.Render("on")
	}

	s += fmt.Sprintf("%s %s", value, item.title)

	return s
}

func (item *formItem) numberItemView() string {
	if item.value <= 0 {
		return fmt.Sprintf("%s %s", offStyle.Render("None"), item.title)
	}

	return fmt.Sprintf("%v %s", item.value, item.title)
}
