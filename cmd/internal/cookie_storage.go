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
			Name: cookieName,
			// Get the first step
			Value:    recipes.NextStep(-1).String(),
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
