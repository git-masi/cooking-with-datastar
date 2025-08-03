package main

import (
	"bytes"
	"cooking-with-datastar/cmd/internal"
	"cooking-with-datastar/cmd/recipes"
	"cooking-with-datastar/cmd/view/about"
	"cooking-with-datastar/cmd/view/cooking"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

var prepTime = map[string]time.Duration{
	"cook-the-chicken": 25 * time.Second,
}

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

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(
			cooking.Recipe(recipe, step, map[string]bool{}),
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var gathered map[string]bool
		err = json.Unmarshal(body, &gathered)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for _, v := range recipe.ListIngredients() {
			if _, ok := gathered[v.Key]; !ok {
				logger.Error("Invalid ingredient", slog.String("key", v.Key))
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}

		cs := internal.NewCookieStorage(recipe, w, r)

		cookie, err := cs.GetIngredientsCookie()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		cookie.Value = hex.EncodeToString(body)

		http.SetCookie(w, cookie)
	})

	mux.HandleFunc("GET /prep/{recipe}/{task}", func(w http.ResponseWriter, r *http.Request) {
		// TODO: sanitize input
		recipe := r.PathValue("recipe")
		task := r.PathValue("task")

		logger.Debug("Prep task", slog.String("recipe", recipe), slog.String("task", task))

		duration, ok := prepTime[task]
		if !ok {
			return
		}

		count := int(duration.Seconds())
		sse := datastar.NewSSE(w, r)

		buf := bytes.NewBuffer([]byte{})
		cooking.Timer(task, internal.DisplayMinutesSeconds(count)).Render(r.Context(), buf)
		sse.PatchElements(
			buf.String(),
			datastar.WithSelectorID(fmt.Sprintf("button-%s", task)),
			datastar.WithModeAfter(),
		)

		ticker := time.NewTicker(1 * time.Second)
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					count--

					buf.Reset()
					cooking.Timer(task, internal.DisplayMinutesSeconds(count)).Render(r.Context(), buf)
					sse.PatchElements(buf.String())
				}
			}
		}()

		t := time.NewTimer(duration)
		<-t.C
		ticker.Stop()
		done <- true

		sse.ExecuteScript(`document.querySelector("#ring").remove()`)
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
