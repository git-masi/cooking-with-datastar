package recipes

import (
	"errors"
	"time"
)

type Ingredient struct {
	Key         string
	Description string
}

type Task struct {
	Key         string
	Description string
}

type CookingMethod struct {
	Name     string
	CookTime time.Duration
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

func (r Recipe) ListPrepTasks() []Task {
	switch r {
	case BuffaloChickenDip:
		return []Task{
			{"cook-the-chicken", "Poach your chicken for approximately 25 minutes. When fully cooked, remove from pot and allow to cool until safe to handle."},
			{"shred", "Shred chicken in food processor."},
			{"heat-the-oven", "Preheat the oven to 350 degrees farenheit."},
			{"cube", "Cut the cream cheese into 1 inch cubes."},
			{"warm-the-sauce", "Heat medium sauce pot over medium-low heat. Add the cubed cream cheese, ranch dressing, hot sauce, black pepper, and garlic powder. Whisk constantly until the cream cheese has dissolved. Remove from heat."},
			{"prep-the-pan", "Apply cooking spray to 9x9 inch pan."},
			{"combine", "Combine the shredded chicken, sauce, green onions, and cheese in a large pot. Transfer to baking pan."},
		}

	case ChocolateChipCookies:
		return []Task{
			{"heat-the-oven", "Preheat the oven to 350 degrees farenheit."},
		}

	case PulledPork:
		return []Task{
			{"combine", "Combine onions, pork, and stock in slow cooker."},
		}

	default:
		return []Task{}
	}
}

func (r Recipe) GetCookingMethod() CookingMethod {
	switch r {
	case BuffaloChickenDip:
		return CookingMethod{"bake", 25 * time.Second}

	case ChocolateChipCookies:
		return CookingMethod{"bake", 12 * time.Second}

	case PulledPork:
		return CookingMethod{"slow cook", 30 * time.Second}

	default:
		return CookingMethod{}
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
		if v.Key == key {
			return v, nil
		}
	}

	return Task{}, errors.New("invalid task")
}
