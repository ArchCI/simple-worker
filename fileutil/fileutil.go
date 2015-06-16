package fileutil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"syscall"
)

func NonblockReadFile(filename string) {
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
