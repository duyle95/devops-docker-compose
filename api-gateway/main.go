package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type GlobalState string

const (
	INIT     GlobalState = "INIT"
	PAUSED   GlobalState = "PAUSED"
	RUNNING  GlobalState = "RUNNING"
	SHUTDOWN GlobalState = "SHUTDOWN"
)

var globalState GlobalState = INIT

var nextStateMapper = map[GlobalState][]GlobalState{
	INIT:     {INIT, RUNNING},
	RUNNING:  {INIT, RUNNING, PAUSED, SHUTDOWN},
	PAUSED:   {INIT, PAUSED, RUNNING, SHUTDOWN},
	SHUTDOWN: {},
}

var logMutex sync.Mutex

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.43")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Put("/state", handleChangeState)

	r.With(StateMiddleware).Get("/state", getCurrentState)
	r.With(StateMiddleware).Get("/request", getContainersInfo)
	r.With(StateMiddleware).Get("/run-log", getRunLog)

	fmt.Println("api-gateway is running on port 8197")

	http.ListenAndServe(":8197", r)
}

func StateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if globalState == PAUSED || globalState == SHUTDOWN {
			http.Error(w, "Service is unavailable", http.StatusServiceUnavailable)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getCurrentState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Accept", "text/plain")
	w.Write([]byte(globalState))
}

func logStateChangeToFile(oldState GlobalState, newState GlobalState) {
	f, err := os.OpenFile("state-changelog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	logMutex.Lock()
	defer logMutex.Unlock()

	logger := log.New(f, "", 0)
	now := time.Now().UTC()
	formattedTime := now.Format("2006-01-02T15.04.05.000Z")
	logger.Printf("%s: %s -> %s\n", formattedTime, oldState, newState)
}

func handleChangeState(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	nextState := GlobalState(body)

	if !slices.Contains(nextStateMapper[globalState], nextState) {
		http.Error(w, "Invalid state transition", http.StatusBadRequest)
		return
	}

	if globalState == INIT && nextState == RUNNING {
		username, password, ok := r.BasicAuth()
		fmt.Printf("username [%s], password [%s] \n", username, password)

		authUsername := getEnv("AUTH_USERNAME", "admin")
		authPassword := getEnv("AUTH_PASSWORD", "admin")
		if !ok || username != authUsername || password != authPassword {
			http.Error(w, "Unauthorized to change state from INIT to RUNNING", http.StatusUnauthorized)
			return
		}
	}

	if nextState == SHUTDOWN {
		shutdownAllContainers()
	}

	fmt.Printf("Changing from %s to %s \n", globalState, nextState)
	oldState := globalState
	globalState = nextState
	logStateChangeToFile(oldState, nextState)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Accept", "text/plain")
	w.Write([]byte("OK"))
}

func shutdownAllContainers() {
	res, err := http.Post("http://node-service:3000/api/stop-all-containers", "", nil)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("error stopping all containers: %s\n", res.Status)
	}
}

func getContainersInfo(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://node-service:3000/api/get-container-info")
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition",
		"attachment; filename=duyle-container-info.txt")
	w.Write(resBody)
}

func getRunLog(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("state-changelog.txt")
	if err != nil {
		fmt.Printf("error reading file: %s\n", err)
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Accept", "text/plain")
	w.Write(content)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
