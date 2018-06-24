package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

// this file implements a simple callback handler for OAuth2 flows

type callbackPayload struct {
	code  string
	state string
}

const (
	responseTimeout = 5
)

var ErrTimeout = errors.New("Timeout waiting for OAuth login")

// makes a func compatible with HandleFunc that pushes the contents of
// a callback into the passed channel
func makeHandler(c chan callbackPayload) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			log.Printf("Could not parse callback req: %v", err)
			http.Error(w, "Bad payload", 400)
			return
		}

		state := req.Form.Get("state")
		code := req.Form.Get("code")

		if (state == "") || (code == "") {
			log.Print("Callback missing params code & state")
			http.Error(w, "Bad request", 400)
			return
		}

		c <- callbackPayload{code: code, state: state}

		w.WriteHeader(200)
		w.Write([]byte("You may close this page"))
	}
}

// Spawns a server listening on `addr` for a OAuth callback with state and code
// set as url parameters to the endpoint. Exits after timeout or on receipt of
// a code/state pair.
func waitForCallback(addr, ep string, timeoutSeconds int) (code, state string, err error) {
	c := make(chan callbackPayload)

	mux := http.NewServeMux()
	mux.Handle(ep, http.HandlerFunc(makeHandler(c)))

	srv := &http.Server{Addr: addr, Handler: mux}

	// start the server in a goroutine so we can shut it down from the main
	// goroutine
	go func() {
		log.Printf("Listening on %v for OAuth callback on %v", addr, ep)
		if err := srv.ListenAndServe(); (err != nil) && (err != http.ErrServerClosed) {
			log.Fatalf("OAuth callback server error: %v", err)
		}
	}()

	var result callbackPayload

	select {
	case result = <-c:
		// we got a state, code pair from the handler
		state, code = result.state, result.code
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		// timed out waiting for callback
		err = ErrTimeout
	}

	// shutdown the server
	log.Printf("Shutting down server on %v", addr)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	srv.Shutdown(ctx)
	return
}
