package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type GlobalState string

const (
	INIT     GlobalState = "INIT"
	PAUSED   GlobalState = "PAUSED"
	RUNNING  GlobalState = "RUNNING"
	SHUTDOWN GlobalState = "SHUTDOWN"
)

var globalState GlobalState = INIT

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.43")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/request", getContainersInfo)

	fmt.Println("api-gateway is running on port 8197")

	http.ListenAndServe(":8197", r)
}

func getContainersInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
