package internal

import "fmt"

func DisplayMinutesSeconds(seconds int) string {
	if seconds < 0 {
		return "00:00"
	}

	if seconds > 3599 {
		return "59:59"
	}

	minutes := seconds / 60
	_seconds := seconds % 60

	return fmt.Sprintf(
		"%s:%s",
		Ternary(minutes < 10, fmt.Sprintf("0%d", minutes), fmt.Sprint(minutes)),
		Ternary(_seconds < 10, fmt.Sprintf("0%d", _seconds), fmt.Sprint(_seconds)),
	)
}
