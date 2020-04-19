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

	var taskRegistry map[string]tasks.TaskType = make(map[string]tasks.TaskType, 0)
	taskRegistry["GOL"] = gameoflife.GOLTaskType

	messenger.InitializeMessenger()

	wm := workers.GetWorkerManager()
	wm.Start()

	_ = tasks.Start(taskRegistry, wm.MapTaskQueue)

	r := RegisterRoutes()
	log.Fatal(http.ListenAndServe(":2216", r))
}
