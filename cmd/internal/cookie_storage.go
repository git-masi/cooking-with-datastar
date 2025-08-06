package internal

import (
	"cooking-with-datastar/cmd/recipes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type CookieStorage struct {
	recipe recipes.Recipe
	res    http.ResponseWriter
	req    *http.Request
}

func NewCookieStorage(recipe recipes.Recipe, res http.ResponseWriter, req *http.Request) CookieStorage {
	return CookieStorage{
		recipe,
		res,
		req,
	}
}

func (cs CookieStorage) GetStepCookie() (*http.Cookie, error) {
	recipeName := cs.recipe.String()
	cookieName := recipeName + "-step"

	cookie, err := cs.req.Cookie(cookieName)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			return nil, err
		}

		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    recipes.GetFirstStep().String(),
			Path:     "/",
			MaxAge:   int((24 * time.Hour).Seconds()),
			HttpOnly: true,                 // Do not allow JS to modify the cookie
			Secure:   true,                 // Only use HTTPS (and localhost)
			SameSite: http.SameSiteLaxMode, // Send cookie when navigating *to* our site
		}

		http.SetCookie(cs.res, cookie)
	}

	return cookie, nil
}

func (cs CookieStorage) ToNextStep() (*http.Cookie, error) {
	cookie, err := cs.GetStepCookie()
	if err != nil {
		return nil, err
	}

	step, err := recipes.ParseRecipeStep(cookie.Value)
	if err != nil {
		return nil, err
	}

	cookie.Path = "/"
	cookie.Value = step.GetNextStep().String()

	return cookie, nil
}

func (cs CookieStorage) GetIngredientsCookie() (*http.Cookie, error) {
	recipeName := cs.recipe.String()
	cookieName := recipeName + "-ingredients"

	cookie, err := cs.req.Cookie(cookieName)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			return nil, err
		}

		gathered := map[string]bool{}

		for _, v := range cs.recipe.ListIngredients() {
			gathered[v.Key] = false
		}

		json, err := json.Marshal(gathered)
		if err != nil {
			return nil, err
		}

		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    hex.EncodeToString(json),
			Path:     "/",
			MaxAge:   int((24 * time.Hour).Seconds()),
			HttpOnly: true,                 // Do not allow JS to modify the cookie
			Secure:   true,                 // Only use HTTPS (and localhost)
			SameSite: http.SameSiteLaxMode, // Send cookie when navigating *to* our site
		}

		http.SetCookie(cs.res, cookie)
	}

	return cookie, nil
}

func (cs CookieStorage) GetPrepCookie(task recipes.Task) (*http.Cookie, error) {
	recipeName := cs.recipe.String()
	cookieName := recipeName + "-prep-" + task.Key

	cookie, err := cs.req.Cookie(cookieName)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			return nil, err
		}

		cookie = &http.Cookie{
			Name:     cookieName,
			Value:    task.PrepTime.String(),
			Path:     "/",
			MaxAge:   int((24 * time.Hour).Seconds()),
			HttpOnly: true,                 // Do not allow JS to modify the cookie
			Secure:   true,                 // Only use HTTPS (and localhost)
			SameSite: http.SameSiteLaxMode, // Send cookie when navigating *to* our site
		}

		http.SetCookie(cs.res, cookie)
	}

	return cookie, nil
}

func (cs CookieStorage) DecrementPrepCookie(task recipes.Task, amount time.Duration) (*http.Cookie, error) {
	cookie, err := cs.GetPrepCookie(task)
	if err != nil {
		return nil, err
	}

	timeRemaining, err := time.ParseDuration(cookie.Value)
	if err != nil {
		return nil, err
	}

	if timeRemaining.Seconds() == 0 {
		return cookie, err
	}

	cookie.Path = "/"
	cookie.Value = (timeRemaining - amount).String()

	http.SetCookie(cs.res, cookie)

	return cookie, nil
}

func (cs CookieStorage) IsEveryTaskFinished() bool {
	for _, t := range cs.recipe.ListPrepTasks() {
		c, _ := cs.GetPrepCookie(t)
		v, _ := time.ParseDuration(c.Value)
		if v.Seconds() != 0 {
			return false
		}
	}

	return true
}
