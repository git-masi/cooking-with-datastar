package main

import (
	"bytes"
	"cooking-with-datastar/cmd/internal"
	"cooking-with-datastar/cmd/view/about"
	"cooking-with-datastar/cmd/view/cooking"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

type BuffaloChickenIngredients struct {
	Chicken          bool `json:"chicken"`
	CreamCheese      bool `json:"cream-cheese"`
	RanchDressing    bool `json:"ranch-dressing"`
	HotSauce         bool `json:"hot-sauce"`
	BlackPepper      bool `json:"black-pepper"`
	GarlicPowder     bool `json:"garlic-powder"`
	GreenOnion       bool `json:"green-onion"`
	MozzarellaCheese bool `json:"mozzarella-cheese"`
	CheddarCheese    bool `json:"cheddar-cheese"`
}

var prepTime = map[string]time.Duration{
	"cook-the-chicken": 25 * time.Second,
}

func main() {
	port := flag.Int("port", 8080, "A port to listen on")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /recipe/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe := r.PathValue("recipe")

		w.Header().Add("content-type", "text/html")

		cooking.Recipe(recipe).Render(r.Context(), w)
	})

	mux.HandleFunc("PATCH /ingredients/{recipe}", func(w http.ResponseWriter, r *http.Request) {
		recipe := r.PathValue("recipe")

		ingredients := &BuffaloChickenIngredients{}
		if err := datastar.ReadSignals(r, ingredients); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Info("PATCH ingredients", slog.String("recipe", recipe), slog.Any("ingredients", ingredients))
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
