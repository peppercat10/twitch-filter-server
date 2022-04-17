package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

func getActiveConsoles() []string {
	files, err := ioutil.ReadDir("./twitch-platform-filter/html")
	if err != nil {
		return nil
	}

	var activeConsoles []string
	for _, file := range files {
		filename := file.Name()
		if strings.Contains(filename, ".html") {
			filenameSplit := strings.Split(filename, ".")
			_, err := strconv.Atoi(filenameSplit[0])
			if err == nil {
				activeConsoles = append(activeConsoles, filenameSplit[0])
			}
		}
	}
	return activeConsoles
}

func refreshLiveGames() {
	fmt.Println("Refreshing live games.")
	activeConsoles := getActiveConsoles()
	if activeConsoles == nil {
		fmt.Println("No active consoles found!")
		return
	}

	commands := append([]string{"twitch-platform-filter/twitch-filter.py", "silent|"}, activeConsoles...)

	cmd := exec.Command("python", commands...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func main() {
	fmt.Println("Executing.")
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hours().Do(refreshLiveGames)
	s.StartAsync()
	port := ":" + os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	http.ListenAndServe(port, http.FileServer(http.Dir("./twitch-platform-filter/html")))
}
