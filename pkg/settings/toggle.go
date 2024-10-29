package settings

type toggle struct {
	id    string
	title string
	value bool
}

func (t *toggle) Toggle() {
	t.value = !t.value
}

func (t *toggle) View() string {
	s := ""

	if t.value {
		s += "ON"
	} else {
		s += "OFF"
	}

	s += " " + t.title

	return s
}
