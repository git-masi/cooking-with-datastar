package main

import (
	"cooking-with-datastar/cmd/internal"
	"cooking-with-datastar/cmd/recipes"
	"cooking-with-datastar/cmd/view/cooking"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

//go:embed "static"
var Files embed.FS

func main() {
	port := flag.Int("port", 8080, "A port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(Files))

	mux.HandleFunc("GET /recipe/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		cs := internal.NewCookieStorage(recipe, r)

		cookie, err := cs.GetStepCookie()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, cookie)

		step, err := recipes.ParseRecipeStep(cookie.Value)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		gathered, err := cs.GetGatheredIngredients()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		finishedTasks, err := cs.GetFinishedTasks()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		finishedCooking, err := cs.FinishedCooking()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(
			cooking.Recipe(recipe, step, gathered, finishedTasks, finishedCooking),
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

		cs := internal.NewCookieStorage(recipe, r)

		cookie, err := cs.GatherIngredients(r.Form)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, cookie)

		finished, err := cs.FinishedGatheringIngredients(cookie)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !finished {
			return
		}

		cookie, err = cs.ToNextStep()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/recipe/"+recipe.String(), http.StatusSeeOther)
	})

	mux.HandleFunc("PATCH /prep/{recipe}/{task}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			logger.Error("Cannot parse recipe", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		task, err := recipes.ParseTask(recipe, r.PathValue("task"))
		if err != nil {
			logger.Error("Cannot parse task", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cs := internal.NewCookieStorage(recipe, r)

		cookie, err := cs.FinishTask(task)
		if err != nil {
			logger.Error("Cannot parse task", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, cookie)
		r.AddCookie(cookie)

		finished, err := cs.FinishedAllTasks()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !finished {
			return
		}

		cookie, err = cs.ToNextStep()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/recipe/"+recipe.String(), http.StatusSeeOther)
	})

	mux.HandleFunc("GET /cook/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		cs := internal.NewCookieStorage(recipe, r)

		cookie, err := cs.GetCookingMethodCookie()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, cookie)

		timeRemaining, err := cs.GetRemainingCookTime()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)
		id := "count-down-" + recipe.String()
		path := fmt.Sprintf("/cook/%s", recipe.String())

		if timeRemaining.Seconds() <= 0 {
			sse.PatchElementTempl(
				cooking.Timer(id, path, 0),
				datastar.WithSelectorID("button-"+recipe.GetCookingMethod().Name),
				datastar.WithModeAfter(),
			)

			sse.ExecuteScript(`document.querySelector("#ring").remove()`)
			return
		}

		seconds := int(timeRemaining.Seconds())

		err = sse.PatchElementTempl(
			cooking.Timer(id, path, seconds),
			datastar.WithSelectorID("button-"+recipe.GetCookingMethod().Name),
			datastar.WithModeAfter(),
		)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ticker := time.NewTicker(1 * time.Second)
		timer := time.NewTimer(timeRemaining)
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					seconds--

					err := sse.PatchElementTempl(
						cooking.Timer(id, path, seconds),
						datastar.WithModeReplace(),
					)
					if err != nil {
						logger.Error(err.Error())
						return
					}
				}
			}
		}()

		<-timer.C
		ticker.Stop()
		done <- true

		sse.ExecuteScript(`document.querySelector("#ring").remove()`)

		sse.PatchElements(
			fmt.Sprintf(`
				<img id="finished-recipe" src="%s"/>
				<button id="button-%s" disabled>%s</button>
			`,
				recipe.GetImageSrc(),
				recipe.GetCookingMethod().Name,
				recipe.GetCookingMethod().Name,
			),
		)
	})

	mux.HandleFunc("PATCH /cook/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe, err := recipes.ParseRecipe(r.PathValue("recipe"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		cs := internal.NewCookieStorage(recipe, r)

		cookie, err := cs.DecrementCookingMethodCookie(1 * time.Second)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, cookie)
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
