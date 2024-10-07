package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.43")
	http.HandleFunc("/get-container-info", getContainerInfoFile)

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getContainerInfoFile(w http.ResponseWriter, r *http.Request) {
	fileName := "containers_info.txt"
	outputFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	writeContainerDetailsToFile(outputFile)

	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error reading file.", http.StatusInternalServerError)
	}
}

func writeContainerDetailsToFile(outputFile *os.File) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	expectedContainerName := []string{"node-service", "golang-service"}

	for _, container := range containers {
		if !strings.Contains(container.Image, expectedContainerName[0]) && !strings.Contains(container.Image, expectedContainerName[1]) {
			continue
		}

		writeOutputToFile(outputFile, fmt.Sprintf("Service: %s\n\n", container.Image))

		getIPAddressAndUptimeAndWriteToFile(dockerClient, container.ID, outputFile)

		execCommandAndWriteOutputToFile(dockerClient, container.ID, outputFile, "List of running processes:\n", "ps", "-a")

		execCommandAndWriteOutputToFile(dockerClient, container.ID, outputFile, "Available disk space:\n", "df", "-h")

		writeOutputToFile(outputFile, "\n\n")
	}
}

func getIPAddressAndUptimeAndWriteToFile(cli *client.Client, containerID string, file *os.File) {
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Printf("Error when inspecting container: %v", err)
		return
	}

	if containerInfo.NetworkSettings != nil {
		for _, netDetails := range containerInfo.NetworkSettings.Networks {
			writeOutputToFile(file, fmt.Sprintf("IP Address: %s\n\n", netDetails.IPAddress))
		}
	}

	startedAt, err := time.Parse(time.RFC3339, containerInfo.State.StartedAt)
	if err != nil {
		log.Printf("Error when parse container started at time: %v", err)
		return
	}

	writeOutputToFile(file, fmt.Sprintf("Time since last boot: %f minutes\n\n", time.Since(startedAt).Minutes()))
}

func execCommandAndWriteOutputToFile(cli *client.Client, containerID string, file *os.File, cmdDesc string, cmd ...string) {
	respID, err := cli.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		Tty:          true,
	})
	if err != nil {
		log.Printf("Error when creating exec instance in container: %v", err)
		return
	}

	resp, err := cli.ContainerExecAttach(context.Background(), respID.ID, container.ExecAttachOptions{})
	if err != nil {
		log.Printf("Error when attaching to exec instance in container: %v", err)
		return
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		log.Printf("Error when reading exec output from container: %v", err)
		return
	}
	writeOutputToFile(file, fmt.Sprintf("%s %s\n", cmdDesc, string(output)))
}

func writeOutputToFile(file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}
