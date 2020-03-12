package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("hello")

	wm := GetWorkerManager()
	uuid := wm.AddWorker(net.IPv4(128, 0, 0, 1))
	wm.RemoveWorker(uuid)
}
