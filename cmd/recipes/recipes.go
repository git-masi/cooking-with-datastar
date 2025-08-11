package recipes

import (
	"errors"
	"time"
)

type Ingredient struct {
	Name        string
	Description string
}

type Task struct {
	Name         string
	Description  string
	Dependencies []string
}

type CookingMethod struct {
	Name        string
	Description string
	CookTime    time.Duration
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
			{"butter", "1 cup butter, softened"},
			{"white-sugar", "1 cup white sugar"},
			{"brow-sugar", "1 cup packed brown sugar"},
			{"eggs", "2 large eggs"},
			{"vanilla", "2 teaspoons vanilla extract"},
			{"baking-soda", "1 teaspoon baking soda"},
			{"hot-water", "2 teaspoons hot water"},
			{"salt", "0.5 teaspoon salt"},
			{"flour", "3 cups all-purpose flour"},
			{"chocolate-chips", "2 cups semisweet chocolate chips"},
			{"walnuts", "1 cup chopped walnuts"},
		}

	case PulledPork:
		return []Ingredient{
			{"pork-shoulder", "3 pound boneless pork shoulder roast"},
			{"ketchup", "1 cup ketchup"},
			{"brown-sugar", "0.5 cup firmly packed brown sugar"},
			{"vinegar", "0.25 cup apple cider vinegar"},
			{"hot-sauce", "Hot sauce to taste"},
		}

	default:
		return []Ingredient{}
	}
}

func (r Recipe) ListPrepTasks() []Task {
	switch r {
	case BuffaloChickenDip:
		return []Task{
			{"cook-the-chicken", "Poach the chicken for approximately 25 minutes. When fully cooked, remove from pot and allow to cool until safe to handle.", []string{}},
			{"shred", "Shred chicken in food processor.", []string{"cook-the-chicken"}},
			{"heat-the-oven", "Preheat the oven to 350 degrees farenheit.", []string{}},
			{"cube", "Cut the cream cheese into 1 inch cubes.", []string{}},
			{"warm-the-sauce", "Heat medium sauce pot over medium-low heat. Add the cubed cream cheese, ranch dressing, hot sauce, black pepper, and garlic powder. Whisk constantly until the cream cheese has dissolved. Remove from heat.", []string{"cube"}},
			{"prep-the-pan", "Apply cooking spray to 9x9 inch pan.", []string{}},
			{"combine", "Combine the shredded chicken, sauce, green onions, and cheese in a large pot. Transfer to baking pan.", []string{"cook-the-chicken", "shred", "heat-the-oven", "cube", "warm-the-sauce", "prep-the-pan"}},
		}

	case ChocolateChipCookies:
		return []Task{
			{"heat-the-oven", "Preheat the oven to 350 degrees farenheit.", []string{}},
			{"beat-eggs", "Beat in eggs, one at a time, then stir in vanilla.", []string{}},
			{"add-baking-soda", "Dissolve baking soda in hot water. Add to batter along with salt.", []string{"beat-eggs"}},
			{"stir-in-flour", "Stir in flour, chocolate chips, and walnuts.", []string{"add-baking-soda"}},
			{"place-dough", "Drop spoonfuls of dough 2 inches apart onto ungreased baking sheets.", []string{"stir-in-flour"}},
		}

	case PulledPork:
		return []Task{
			{"place", "Place pork roast in a slow cooker.", []string{}},
			{"combine", "Whisk ketchup, brown sugar, vinegar, and hot sauce together in a bowl until well combined", []string{}},
			{"pour", "Pour the mixture over the pork. Turn pork to coat completely.", []string{"place", "combine"}},
		}

	default:
		return []Task{}
	}
}

func (r Recipe) GetCookingMethod() CookingMethod {
	switch r {
	case BuffaloChickenDip:
		return CookingMethod{"bake", "Bake for 20-30 minutes, or until the cheese has melted and the sides are starting to bubble.", 10 * time.Second}

	case ChocolateChipCookies:
		return CookingMethod{"bake", "Bake for 10-12 minutes", 5 * time.Second}

	case PulledPork:
		return CookingMethod{"slow cook", "Slow cook on low for 8 to 10 hours or High for 4 to 6 hours.", 15 * time.Second}

	default:
		return CookingMethod{}
	}
}

func (r Recipe) GetImageSrc() string {
	switch r {
	case BuffaloChickenDip:
		return "/static/buffalo_chicken_dip_pixel_art_small.png"

	case ChocolateChipCookies:
		return ""

	case PulledPork:
		return ""

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

func (s Step) String() string {
	return stepName[s]
}

func (s Step) GetNextStep() Step {
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

func GetFirstStep() Step {
	return Gather
}

func ParseRecipeStep(name string) (Step, error) {
	for r, n := range stepName {
		if n == name {
			return r, nil
		}
	}

	return -1, errors.New("invalid recipe name")
}

func ParseTask(r Recipe, key string) (Task, error) {
	for _, v := range r.ListPrepTasks() {
		if v.Name == key {
			return v, nil
		}
	}

	return Task{}, errors.New("invalid task")
}
