package recipes

import "errors"

type Ingredient struct {
	Key         string
	Description string
}

type Recipe int

const (
	BuffaloChickenDip Recipe = iota
	ChocolateChipCookies
	PulledPork
)

var recipeName = map[Recipe]string{
	BuffaloChickenDip:    "buffalo-chicken-dip",
	ChocolateChipCookies: "chocolate-chip-cookies",
	PulledPork:           "pulled-pork",
}

func (r Recipe) String() string {
	return recipeName[r]
}

func (r Recipe) ListIngredients() []Ingredient {
	switch r {
	case BuffaloChickenDip:
		return []Ingredient{
			{"chicken", "3 large boneless skinless chicken breasts"},
			{"cream-cheese", "8 ounces cream cheese"},
			{"ranch-dressing", "1 cup ranch dressing"},
			{"hot-sauce", "1 cup hot sauce"},
			{"black-pepper", "1 teaspoon freshly ground black pepper"},
			{"garlic-powder", "1 teaspoon garlic powder"},
			{"green-onion", "0.5 cup green onion"},
			{"mozzarella-cheese", "1.5 cups mozzarella cheese"},
			{"cheddar-cheese", "1.5 cups cheddar cheese"},
		}

	case ChocolateChipCookies:
		return []Ingredient{
			{"chocolate-chips", "1 cup chocolate chips"},
		}

	case PulledPork:
		return []Ingredient{
			{"pork-shoulder", "3 pounds pork shoulder"},
		}

	default:
		return []Ingredient{}
	}
}
func (r Recipe) GetCookingMethod() string {
	switch r {
	case BuffaloChickenDip:
		return "bake"

	case ChocolateChipCookies:
		return "bake"

	case PulledPork:
		return "slow cook"

	default:
		return ""
	}
}

func ListRecipes() []Recipe {
	list := []Recipe{}

	for r := range recipeName {
		list = append(list, r)
	}

	return list
}

func ParseRecipe(name string) (Recipe, error) {
	for r, n := range recipeName {
		if n == name {
			return r, nil
		}
	}

	return -1, errors.New("invalid recipe name")
}

type Step int

const (
	Gather Step = iota
	Prepare
	Cook
	Done
)

var stepName = map[Step]string{
	Gather:  "gather",
	Prepare: "prepare",
	Cook:    "cook",
}

func (r Step) String() string {
	return stepName[r]
}

func (s Step) Next() Step {
	switch s {
	case Gather:
		return Prepare

	case Prepare:
		return Cook

	case Cook:
		return Done

	default:
		return Gather
	}
}

func ParseRecipeStep(name string) (Step, error) {
	for r, n := range stepName {
		if n == name {
			return r, nil
		}
	}

	return -1, errors.New("invalid recipe name")
}
