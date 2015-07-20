package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/http"

	"gopkg.in/yaml.v2"
	//"github.com/garyburd/redigo/redis"
	log "github.com/Sirupsen/logrus"

	"github.com/ArchCI/simple-worker/fileutil"
	"github.com/ArchCI/simple-worker/dbutil"
)

// The task struct to run test
/*
type Task struct {
	Id      int64 `json:"id"`
	Commit  string  `json:"commit"`
	Public  bool    `json:"is_public"`
	Type    string  `json:"type"`
	Project string  `json:"project`
	Url     string  `json:"url"`
}
*/

// Check if the worker can run task or not
func checkRequirement() bool {
	// TODO: it should not work for mac os but docker may be installed
	_, err := exec.LookPath("docker")
	if err != nil {
		return false
	} else {
		return true
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

// Struct for archci.yml
type ArchciConfig struct {
	Image  string   `yaml:"image"`
	Script []string `yaml:"script"`
}

// Parse archci.yml to struct
func ParseYaml(filename string) ArchciConfig {
	// fmt.Println("Start parse yaml") // TODO: Make it as debug log

	var archciConfig ArchciConfig
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &archciConfig)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Value: %#v\n", config.Script[0])
	return archciConfig
}

// Use archci.yml struct to generate archci.sh
func GenerateArchciShellContent(archciConfig ArchciConfig) string {
	// Add this and user's scripts into archci.sh
	baseShellContent := `#!/bin/bash
set -e
cd /project
`

	archciShellContent := baseShellContent
	for _, script := range archciConfig.Script {
		archciShellContent += script + "\n"
	}

	return archciShellContent
}

func PostString(url string, data string) {
	if err := http.Post(url, strings.NewReader(data)); err != nil {
		log.Fatal("could not post")
	}
}

// The simple worker to pull task and run
func main() {
	fmt.Println("Start ArchCI simple worker")

	log.Info("Start simple-worker")

	build := dbutil.GetBuildToTest();
	fmt.Println("Build the project " + build.ProjectName)

	// TODO: Support get parameter from command-line(server url, interval time and task number)

	// Exit if it doesn't meet the requirement
	if checkRequirement() == false {
		// TODO: Add more info to install requirement
		fmt.Println("Need some requirements, exit now")
		return
	}

	// Loop to pull task and run test
	for {

		fmt.Println("Start while loop")

		// HTTP request to get task array
		// TODO(tobe): change to http://archci.com/tasks?number=1
		// TODO(tobe): remove task to use build
		//task := Task{Id: int64(123), Commit: "commit", Public: true, Type: "github", Project: "test-project", Url: "https://github.com/tobegit3hub/test-project.git"}
		//tasks := []Task{task}


		// If no task, sleep and wait for next
		//if len(tasks) == 0 {
		//	fmt.Println("Sleep 5 seconds then pull task again")
		//	time.Sleep(5 * time.Second)
		//	continue
		//}

		// Get the task and run test
		//task = tasks[0]
		//fmt.Println(task.Project)

		// 1. Clone the code in specified directory
		// TODO: support user defined directory, avoid the name conflict
		// TODO: Support other command than "git clone"
		cloneCmd := exec.Command("git", "clone", build.RepoUrl)

		cloneOut, err := cloneCmd.Output()
		if err != nil {
			// TODO: Don't be so easy to exit
			fmt.Println("Error to run clone command")
			os.Exit(1)
		}
		fmt.Println(string(cloneOut)) // Nothing to output if success
		fmt.Println("Success to clone the code")

		// 2. Parse archci.yaml file for base image and test scripts
		archciConfig := ParseYaml(build.ProjectName + "/archci.yml")

		// fmt.Printf("Value: %#v\n", archciConfig.Image)
		dockerImage := archciConfig.Image
		fmt.Printf("Docker image: %#v\n", dockerImage)

		// 3. Generate archci.sh to "cd", combine user scipts and redirect STDOUT to file, this file should put into user's root directory
		// TODO: Make sure that the archci.sh is not conflict or just rm the user's one
		archciShellContent := GenerateArchciShellContent(archciConfig)
		archciShellFile, err2 := os.OpenFile(build.ProjectName+"/archci.sh", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755) // TODO: Make it a constant

		if err2 != nil {
			panic(err2)
		}
		defer archciShellFile.Close()
		_, err2 = archciShellFile.WriteString(archciShellContent)
		if err2 != nil {
			panic(err2)
		}
		archciShellFile.Sync()
		archciShellFile.Close()
		fmt.Println("Success to create archci.sh")

		// 4. Docker run the base image and put the code into container to test
		// docker run --rm -v $PWD:/project golan:1.4 /project/archci.sh > docker.log 2>&1 ; echo $? > exit_code.log
		cpuLimit := ""    // " -c 2 "
		memoryLimit := "" // " -m 100m "
		dockerCommand := "docker run --rm " + cpuLimit + memoryLimit + "-v $PWD/" + build.ProjectName + ":/project " + dockerImage + " /project/archci.sh > docker.log 2>&1 ; echo $? > exit_code.log"

		dockerCmd := exec.Command("sh", "-c", dockerCommand)
		dockerOut, err := dockerCmd.Output()
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(dockerOut))
		fmt.Println("Success to run " + dockerCommand)

		// 5. Non-block read the log and exit_code file and put them into redis
		//fileutil.NonblockReadFile("docker.log")
		fileutil.WriteFileToRedis(build.Id, "docker.log")
		// PostString("http://127.0.0.1:8080/v1/account", "my log one")

		// 6. Delete the code
		// TODO: make it a function to call
		rmCmd := exec.Command("rm", "-rf", build.ProjectName)
		rmOut, err := rmCmd.Output()
		if err != nil {
			// TODO: Don't be so easy to exit
			os.Exit(1)
		}
		fmt.Println(string(rmOut))
		fmt.Println("Success to delete the code")

		// Sleep for next task
		fmt.Println("Sleep 60 seconds for next task")
		time.Sleep(60 * time.Second)
	}

	fmt.Println("Simple worker exists")

}
