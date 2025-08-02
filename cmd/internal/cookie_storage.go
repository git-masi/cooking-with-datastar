package internal

import (
	"cooking-with-datastar/cmd/recipes"
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
			Value:    recipes.Gather.String(),
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
