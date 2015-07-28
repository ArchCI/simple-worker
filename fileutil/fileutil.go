package fileutil

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/ArchCI/simple-worker/redisutil"
)

// NonblockReadFile will non-block read the file.
func NonblockReadFile(filename string) {
	// TODO(tobe): this is not used and needs test.

	// open input file
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	syscall.SetNonblock(int(fi.Fd()), true)

	// make a read buffer
	r := bufio.NewReader(fi)

	for {
		// read a chunk
		line, _, err := r.ReadLine()

		fmt.Println(line)

		if err == io.EOF {
			// do something here
			fmt.Println("It is end of file")
			break
		} else if err != nil {
			panic(err)
		}

	}

}

// WriteFileToRedis takes string to write to redis.
func WriteFileToRedis(buildId int64, logfile string) bool {

	logs := []string{}

	file, err := os.Open(logfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		logs = append(logs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return redisutil.WriteLogsToRedis(buildId, logs)
}

// ReadFile take the path of file and return its content as string.
func ReadFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	return string(data), err
}
