package internal

import "cooking-with-datastar/cmd/recipes"

func GetBorderStyle(current recipes.Step, target recipes.Step) string {
	if current == target {
		return "border: .25rem solid var(--pico-primary); border-radius: 5px;"
	}

	return "border: .25rem solid var(--pico-color-grey-100); border-radius: 5px;"
}
