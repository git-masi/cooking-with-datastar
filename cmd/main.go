package main

import (
	"cooking-with-datastar/cmd/internal"
	"cooking-with-datastar/cmd/recipes"
	"cooking-with-datastar/cmd/view/about"
	"cooking-with-datastar/cmd/view/cooking"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

func main() {
	port := flag.Int("port", 8080, "A port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /recipe/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		cs := internal.NewCookieStorage(recipe, w, r)

		cookie, err := cs.GetStepCookie()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		step, err := recipes.ParseRecipeStep(cookie.Value)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		cookie, err = cs.GetIngredientsCookie()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data, err := hex.DecodeString(cookie.Value)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var gathered map[string]bool
		err = json.Unmarshal(data, &gathered)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(
			cooking.Recipe(recipe, step, gathered),
		)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("PATCH /gather/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = r.ParseForm()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ingredients := recipe.ListIngredients()
		count := 0

		gathered := map[string]bool{}
		for _, v := range ingredients {
			// Form data only includes the name of the checkbox if it is checked.
			// So if the value exists at all then we know the value is "true"
			checked := r.Form.Has(v.Key)
			gathered[v.Key] = checked
			if checked {
				count++
			}
		}

		cs := internal.NewCookieStorage(recipe, w, r)

		cookie, err := cs.GetIngredientsCookie()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(gathered)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		cookie.Path = "/"
		cookie.Value = hex.EncodeToString(data)

		http.SetCookie(w, cookie)

		if count == len(ingredients) {
			cookie, err := cs.GetStepCookie()
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			step, err := recipes.ParseRecipeStep(cookie.Value)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			cookie.Path = "/"
			cookie.Value = recipes.NextStep(step).String()

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/recipe/"+recipe.String(), http.StatusSeeOther)
		}
	})

	mux.HandleFunc("GET /prep/{recipe}/{task}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		task, err := recipes.ParseTask(recipe, r.PathValue("task"))
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		cs := internal.NewCookieStorage(recipe, w, r)

		cookie, err := cs.GetPrepCookie(task)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		timeRemaining, err := time.ParseDuration(cookie.Value)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)

		if timeRemaining.Seconds() == 0 {
			sse.PatchElementTempl(
				cooking.Timer(recipe, task, 0),
				datastar.WithSelectorID(fmt.Sprintf("button-%s", task.Key)),
				datastar.WithModeAfter(),
			)

			sse.ExecuteScript(`document.querySelector("#ring").remove()`)
			return
		}

		seconds := int(task.PrepTime.Seconds())

		err = sse.PatchElementTempl(
			cooking.Timer(recipe, task, seconds),
			datastar.WithSelectorID(fmt.Sprintf("button-%s", task.Key)),
			datastar.WithModeAfter(),
		)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ticker := time.NewTicker(1 * time.Second)
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					seconds--

					err := sse.PatchElementTempl(
						cooking.Timer(recipe, task, seconds),
						datastar.WithModeReplace(),
					)
					if err != nil {
						logger.Error(err.Error())
						return
					}
				}
			}
		}()

		t := time.NewTimer(task.PrepTime)
		<-t.C
		ticker.Stop()
		done <- true

		sse.ExecuteScript(`document.querySelector("#ring").remove()`)
	})

	mux.HandleFunc("PATCH /prep/{recipe}/{task}", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("patch cookie")
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		task, err := recipes.ParseTask(recipe, r.PathValue("task"))
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cs := internal.NewCookieStorage(recipe, w, r)

		cookie, err := cs.GetPrepCookie(task)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		timeRemaining, err := time.ParseDuration(cookie.Value)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if timeRemaining.Seconds() == 0 {
			return
		}

		cookie.Path = "/"
		cookie.Value = (timeRemaining - 1*time.Second).String()

		http.SetCookie(w, cookie)
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Render static HTML
		about.About().Render(r.Context(), w)
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		cooking.Cooking().Render(r.Context(), w)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	logger.Info("Starting server", slog.Int("port", *port))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Cannot start server", "error", err.Error())
	}
}
