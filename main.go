package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eagraf/synchronizer/tasks"
	"github.com/eagraf/synchronizer/tasks/gameoflife"
)

func main() {
	fmt.Println("Starting synchronizer")

	/*wm := GetWorkerManager()
	uuid := wm.AddWorker(net.IPv4(128, 0, 0, 1), "cloud")
	wm.RemoveWorker(uuid)*/

	/*sampleTask := tasks.Task{
		TaskType: "GOL",
		Config: tasks.TaskConfig{
			NumWorkers: 1,
		},
		Input: "Helloo",
		State: "WOOORLD",
	}*/

	var taskRegistry map[string]tasks.TaskType = make(map[string]tasks.TaskType, 0)
	taskRegistry["GOL"] = gameoflife.GOLTaskType

	newIntent := tasks.Intent{
		IntentType: "setup",
		TaskType:   "GOL",
		TaskUUID:   "Hello",
		Config: tasks.TaskConfig{
			NumWorkers: 4,
		},
		Input: "Chicken Butt",
	}

	ts := tasks.Start(taskRegistry)
	ts.IntentQueue <- &newIntent

	r := RegisterRoutes()
	log.Fatal(http.ListenAndServe(":2216", r))
}
