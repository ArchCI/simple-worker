package main

import (
	"fmt"
	"time"
)

// The task struct to run test
type Task struct {
	Id     float64 `json:"id"`
	Commit string  `json:"commit"`
	Public bool    `json:"is_public"`
	Type   string  `json:"type"`
	Url    string  `json:"url"`
}

// Check if the worker can run task or not
func checkRequirement() bool {

	// TODO: check if docker is installed, now always return true
	if true {
		return true
	} else {
		return false
	}

}

// The simple worker to pull task and run
func main() {

	fmt.Println("Start ArchCI simple worker")

	// TODO: Support get parameter from command-line(server url, interval time and task number)

	// Exit if it doesn't meet the requirement
	if checkRequirement() == false {
		// TODO: Add more info to install requirement
		fmt.Println("Need some requirements, exit now")
		return
	}

	// Loop to pull task and run test
	for {

		// HTTP request to get task array
		// TODO: change to http://archci.com/tasks?number=1
		task := Task{Id: 123, Commit: "commit", Public: true, Type: "github", Url: "https"}
		tasks := []Task{task}

		// If no task, sleep and wait for next
		if len(tasks) == 0 {
			fmt.Println("Sleep 5 seconds then pull task again")
			time.Sleep(5 * time.Second)
			continue
		}

		// Get the task and run test
		task = tasks[0]
		fmt.Println(task.Id)

	}

	fmt.Println("Simple worker exists")

}
