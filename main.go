package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hello")

	/*wm := GetWorkerManager()
	uuid := wm.AddWorker(net.IPv4(128, 0, 0, 1), "cloud")
	wm.RemoveWorker(uuid)*/

	r := RegisterRoutes()
	log.Fatal(http.ListenAndServe(":2216", r))
}
