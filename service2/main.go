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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.43")
	http.HandleFunc("/get-container-info", getContainerInfoFile)

	http.HandleFunc("/stop-all-containers", stopAllContainers)

	fmt.Println("golang-service is running on port 3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func stopAllContainers(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to stop all containers")
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

	var golangContainers, otherContainers []types.Container

	for _, c := range containers {
		if strings.Contains(c.Image, "golang") {
			golangContainers = append(golangContainers, c)
		} else {
			otherContainers = append(otherContainers, c)
		}
	}

	for _, c := range otherContainers {
		stopContainer(dockerClient, c)
	}

	for _, c := range golangContainers {
		stopContainer(dockerClient, c)
	}
}

func stopContainer(dockerClient *client.Client, c types.Container) {
	timeout := 10
	err := dockerClient.ContainerStop(context.Background(), c.ID, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		log.Printf("Error stopping container: %v", err)
	}
	log.Printf("Successfully stopped container [ID:%s] - [Image:%s]", c.ID, c.Image)
}

func getContainerInfoFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get all container info")
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

	for _, container := range containers {
		writeOutputToFile(outputFile, fmt.Sprintf("Service: %s\n\n", container.Image))

		writeIPAddressAndUptimeToFile(dockerClient, container.ID, outputFile)

		output, _ := execCommandInContainer(dockerClient, container.ID, "ps", "-a")
		writeOutputToFile(outputFile, fmt.Sprintf("%s %s\n", "List of running processes:\n", string(output)))

		output, _ = execCommandInContainer(dockerClient, container.ID, "df", "-h")
		writeOutputToFile(outputFile, fmt.Sprintf("%s %s\n", "Available disk space:\n", string(output)))

		writeOutputToFile(outputFile, "\n\n")
	}
}

func writeIPAddressAndUptimeToFile(cli *client.Client, containerID string, file *os.File) {
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

func execCommandInContainer(cli *client.Client, containerID string, cmd ...string) ([]byte, error) {
	respID, err := cli.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		Tty:          true,
	})
	if err != nil {
		log.Printf("Error when creating exec instance in container: %v", err)
		return nil, err
	}

	resp, err := cli.ContainerExecAttach(context.Background(), respID.ID, container.ExecAttachOptions{})
	if err != nil {
		log.Printf("Error when attaching to exec instance in container: %v", err)
		return nil, err
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		log.Printf("Error when reading exec output from container: %v", err)
		return nil, err
	}

	return output, nil
}

func writeOutputToFile(file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}
