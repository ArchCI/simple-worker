package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"syscall"
)

// The task struct to run test
type Task struct {
	Id     float64 `json:"id"`
	Commit string  `json:"commit"`
	Public bool    `json:"is_public"`
	Type   string  `json:"type"`
	Project string `json:"project`
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

// Execute the command to replace current process. TODO: not used yet
func Exec(args []string) {
	path, err := exec.LookPath(args[0])
	if err != nil {
		os.Exit(1)
	}
	err = syscall.Exec(path, args, os.Environ())
	if err != nil {
		os.Exit(1)
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
		task := Task{Id: 123, Commit: "commit", Public: true, Type: "github", Project: "test-project", Url: "https://github.com/tobegit3hub/test-project.git"}
		tasks := []Task{task}

		// If no task, sleep and wait for next
		if len(tasks) == 0 {
			fmt.Println("Sleep 5 seconds then pull task again")
			time.Sleep(5 * time.Second)
			continue
		}

		// Get the task and run test
		task = tasks[0]
		fmt.Println(task.Project)

		// 1. Clone the code in specified directory
		// TODO: support user defined directory, avoid the name conflict
		// TODO: Support other command than "git clone"
		cloneCmd := exec.Command("git", "clone", task.Url)
		cloneOut, err := cloneCmd.Output()
		if err != nil {
			// TODO: Don't be so easy to exit
			os.Exit(1)
		}
		fmt.Println(string(cloneOut)) // Nothing to output if success
		fmt.Println("Success to clone the code")

		// 2. Parse archci.yaml file for base image and test scripts





		// Delete the code
		// TODO: make it a function to call
		rmCmd := exec.Command("rm", "-rf", task.Project)
		rmOut, err := rmCmd.Output()
		if err != nil {
			// TODO: Don't be so easy to exit
			os.Exit(1)
		}
		fmt.Println(string(rmOut))
		fmt.Println("Success to delete the code")


		time.Sleep(100 * time.Second)
	}

	fmt.Println("Simple worker exists")

}
