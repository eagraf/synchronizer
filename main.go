package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eagraf/synchronizer/messenger"

	"github.com/eagraf/synchronizer/tasks"
	"github.com/eagraf/synchronizer/tasks/gameoflife"
	"github.com/eagraf/synchronizer/workers"
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
		Input: map[string]interface{}{
			"size": 8,
			"board": []int8{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 1, 0, 0, 0,
				0, 0, 0, 0, 0, 1, 0, 0,
				0, 0, 0, 1, 1, 1, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
	}

	messenger.InitializeMessenger()

	wm := workers.GetWorkerManager()
	wm.Start()

	ts := tasks.Start(taskRegistry, wm.MapTaskQueue)
	ts.IntentQueue <- &newIntent

	r := RegisterRoutes()
	log.Fatal(http.ListenAndServe(":2216", r))
}
