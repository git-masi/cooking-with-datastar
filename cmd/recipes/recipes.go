package recipes

import "errors"

type Ingredient struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Gathered    bool   `json:"gathered"`
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
			{"chicken", "3 large boneless skinless chicken breasts", false},
			{"cream-cheese", "8 ounces cream cheese", false},
			{"ranch-dressing", "1 cup ranch dressing", false},
			{"hot-sauce", "1 cup hot sauce", false},
			{"black-pepper", "1 teaspoon freshly ground black pepper", false},
			{"garlic-powder", "1 teaspoon garlic powder", false},
			{"green-onion", "0.5 cup green onion", false},
			{"mozzarella-cheese", "1.5 cups mozzarella cheese", false},
			{"cheddar-cheese", "1.5 cups cheddar cheese", false},
		}

	case ChocolateChipCookies:
		return []Ingredient{
			{"chocolate-chips", "1 cup chocolate chips", false},
		}

	case PulledPork:
		return []Ingredient{
			{"pork-shoulder", "3 pounds pork shoulder", false},
		}

	default:
		return []Ingredient{}
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
