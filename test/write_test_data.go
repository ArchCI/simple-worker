package main

import (
       "fmt"

       "github.com/ArchCI/simple-worker/fileutil"
)


func main() {

     fmt.Println("Start to write data")


     fileutil.WriteFileToRedis(int64(123), "docker.log")


     fmt.Println("End")

}