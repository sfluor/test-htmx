package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	components "github.com/sfluor/test-htmx/views/components/buldan"
	"github.com/sfluor/test-htmx/views/model/buldan"
)

func returnError(w http.ResponseWriter, msg string, code int) {
	payload := map[string]interface{}{
		"toast": map[string]interface{}{
			"level":   "error",
			"message": msg,
		},
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[error] Failed to marshal toaster error: %s (original msg: %s)", err, msg)
		return
	}

	w.Header().Add("Hx-Trigger", string(raw))
	http.Error(w, msg, code)
}

func renderSSE(ctx context.Context, comp templ.Component, w http.ResponseWriter) {
	fmt.Fprintf(w, "event: message\n")
	var buf bytes.Buffer
	comp.Render(ctx, &buf)
	fmt.Fprintf(w, "data: %s \n\n", buf.String())
	w.(http.Flusher).Flush()
}

func RegisterBuldan(engine buldan.Engine) {

	http.HandleFunc("GET /buldan", func(w http.ResponseWriter, r *http.Request) {
		name := buldan.GenerateName()
		home := components.Home(name)
		log.Print("Get home")
		home.Render(r.Context(), w)
	})

	http.HandleFunc("POST /buldan/new", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			returnError(w, "A name must be provided for the admin of the game", 400)
			return
		}

		game, err := engine.NewGame(name)
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		}

		log.Printf("Creating new game: %+v, admin: %s", game, name)

		w.Header().Add(
			"HX-Redirect",
			fmt.Sprintf("/buldan/instance/%s/%s", game.ID, name),
		)
		w.WriteHeader(301)
	})

	http.HandleFunc("GET /buldan/instance/{id}/{playerID}", func(w http.ResponseWriter, r *http.Request) {
		gameID := r.PathValue("id")
		playerID := r.PathValue("playerID")

		if playerID == "" {
			returnError(w, "Expected non empty player ID", 400)
			return
		}

		ctx := r.Context()

		acceptHeader := strings.ToLower(r.Header.Get("Accept"))
		if !strings.Contains(acceptHeader, "text/event-stream") {
			log.Printf(
				"Accept header was %s returning pre-lobby HTML for (%s/%s): %s",
				acceptHeader,
				gameID,
				playerID,
				r.URL.RequestURI(),
			)
			components.PreLobby(r.URL.RequestURI()).Render(ctx, w)
			return
		}

		game, ch, err := engine.JoinGame(gameID, playerID)
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(200)

		log.Printf("Opening stream for events for (%s/%s)", gameID, playerID)

		renderSSE(ctx, components.Lobby(playerID, game), w)

		for {
			select {
			case evt, ok := <-ch:
				if !ok {
					log.Printf("Game %s for %s is closed", gameID, playerID)
					return
				}

				log.Printf("Game %s for %s, received: %+v", gameID, playerID, evt)
			case <-ctx.Done():
				log.Printf("Player %s in game %s walked away", playerID, gameID)
				return
			}
		}

	})
}
