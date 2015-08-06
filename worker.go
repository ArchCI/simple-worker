package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/http"
	"gopkg.in/yaml.v2"

	"github.com/ArchCI/archci/models"
	"github.com/ArchCI/simple-worker/config"
	"github.com/ArchCI/simple-worker/dbutil"
	"github.com/ArchCI/simple-worker/fileutil"
	"github.com/ArchCI/simple-worker/iputil"
)

// Build with -ldflags "-X main.GitVersion `git rev-parse HEAD` -X main.BuildTime `date -u '+%Y-%m-%d_%I:%M:%S'`"
var (
	GitVersion = "No git version provided"
	BuildTime  = "No build time provided"
)

// CheckRequirement checks if it can run task or not.
func CheckRequirement() bool {
	// TODO: it should not work for mac os but docker may be installed
	_, err := exec.LookPath("docker")
	if err != nil {
		return false
	} else {
		return true
	}
}

// ParseArchciYaml parses .archci.yml to struct.
func ParseArchciYaml(filename string) config.ArchciConfig {
	log.Debug("Start to parse yaml")

	var archciConfig config.ArchciConfig
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &archciConfig)
	if err != nil {
		panic(err)
	}

	return archciConfig
}

// ParseWorkerYaml parses worker.yml to struct.
func ParseWorkerYaml(filename string) config.WorkerConfig {
	log.Debug("Start to parse yaml")

	var workerConfig config.WorkerConfig
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &workerConfig)
	if err != nil {
		panic(err)
	}

	return workerConfig
}

// GenerateArchciShellContent uses archci.yml struct to generate archci.sh.
func GenerateArchciShellContent(archciConfig config.ArchciConfig) string {
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

// Main function to get build and run test.
func main() {

	// Build with git version and build time. Print them when it starts.
	fmt.Println("Git version: " + GitVersion)
	fmt.Println("Build time: " + BuildTime)

	fmt.Println("Start ArchCI simple worker")

	dbutil.InitializeModels()

	log.Info("Start simple-worker")

	workerConfig := ParseWorkerYaml("worker.yml")

	// Record the worker in database.
	ip, _ := iputil.GetLocalIp()
	dbutil.AddWorker(rand.Int63(), ip, time.Now(), models.WORKER_STATUS_BUSY)

	// Exit if it doesn't meet the requirement.
	if CheckRequirement() == false {
		// TODO: Add more info to install requirement
		fmt.Println("Need some requirements, exit now")
		return
	}

	// Loop to pull task and run test.
	for {

		fmt.Println("Start while loop")

		build, err := dbutil.GetOneNotStartBuild()
		if err != nil {
			fmt.Println("No build to run test")

			fmt.Println("Sleep for next task")
			time.Sleep(time.Duration(workerConfig.Interval) * time.Second)
			continue
		}

		//build = models.Build{Id:1234, ProjectName: "test-project", RepoUrl: "https://github.com/tobegit3hub/test-project", Branch: "master"}
		fmt.Println("Build the project " + build.ProjectName)

		dbutil.UpdateBuildStatus(build.Id, models.BUILD_STATUS_BUILDING)

		// 1. Clone the code in specified directory
		// TODO: support user defined directory, avoid the name conflict
		// TODO: Support other command than "git clone"
		cloneCmd := exec.Command("git", "clone", build.RepoUrl)

		cloneOut, err := cloneCmd.Output()
		if err != nil {
			// TODO: Don't be so easy to exit
			fmt.Println("Error to run clone command")
			fmt.Println(string(cloneOut))
			return
		}
		fmt.Println("Success to clone the code")

		// 2. Parse archci.yaml file for base image and test scripts
		archciConfig := ParseArchciYaml(build.ProjectName + "/.archci.yml")
		dockerImage := archciConfig.Image

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
			fmt.Println(string(dockerOut))
			os.Exit(1)
		}
		fmt.Println("Success to run " + dockerCommand)

		// 5. Delete the code
		// TODO: make it a function to call
		rmCmd := exec.Command("rm", "-rf", build.ProjectName)
		rmOut, err := rmCmd.Output()
		if err != nil {
			// TODO: Don't be so easy to exit
			fmt.Println(string(rmOut))
			os.Exit(1)
		}
		fmt.Println("Success to delete the code")

		// 6. Non-block read the log and exit_code file and put them into redis
		fileutil.WriteFileToRedis(build.Id, "docker.log")

		// 7. Read exit code file and update build status
		exitCodeFileContent, err9 := fileutil.ReadFile("exit_code.log")
		if err9 != nil {
			fmt.Println(err9)
			panic(err9)
		}
		exitCode, _ := strconv.Atoi(strings.TrimSpace(exitCodeFileContent))
		if exitCode == 0 {
			fmt.Println("Exit code is 0")
			dbutil.UpdateBuildStatus(build.Id, models.BUILD_STATUS_SUCCESS)
		} else {
			fmt.Println("Exit code is not 0")
			dbutil.UpdateBuildStatus(build.Id, models.BUILD_STATUS_FAIL)
		}

		// 9. Send POST to webhook
		if exitCode == 0 {
			for _, url := range archciConfig.Webhook.Success {
				log.Debug("Trigger webhook to send POST to " + url)

				err = http.Post(url, strings.NewReader("{build: success}"))
				if err != nil {
					log.Fatal("Could not send post request")
					fmt.Println(err)
				}
			}
		} else {
			for _, url := range archciConfig.Webhook.Failure {
				log.Debug("Trigger webhook to send POST to " + url)

				err = http.Post(url, strings.NewReader("{build: failure}"))
				if err != nil {
					log.Fatal("Could not send post request")
					fmt.Println(err)
				}
			}
		}

		// Sleep for next task
		fmt.Println("Sleep for next task")
		time.Sleep(time.Duration(workerConfig.Interval) * time.Second)
	}

	fmt.Println("Simple worker exists")

}
