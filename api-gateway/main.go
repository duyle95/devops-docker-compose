package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.43")
	http.HandleFunc("/api/get-container-info", getContainersInfo)

	fmt.Println("api-gateway is running on port 8197")
	err := http.ListenAndServe(":8197", nil)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
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
