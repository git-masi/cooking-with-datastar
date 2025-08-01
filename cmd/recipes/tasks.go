package recipes

type Task struct {
	Key         string
	Description string
}

func ListPrepTasks(r Recipe) []Task {
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
