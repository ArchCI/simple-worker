package main

import (
	"fmt"

	"github.com/ArchCI/simple-worker/iputil"
)

func main() {
	ip, err := iputil.GetLocalIp()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)
}
