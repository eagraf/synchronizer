package gameoflife

import (
	"math/rand"

	"github.com/eagraf/synchronizer/tasks"
)

// GOLPartialResult represents the results of one workers execution
type GOLPartialResult struct {
	Data string
}

// GOLState represents the current state of the Game of Life task
type GOLState struct {
	round int
}

// GOLTaskType defines the code needed to execute game of life
var GOLTaskType = tasks.TaskType{
	Name:   "GOL",
	Setup:  setup,
	Reduce: reduce,
}

// Setup the game of life task
func setup(intent *tasks.Intent) (*tasks.TaskInstance, []*tasks.Intent) {

	// Initialize partial results and intents
	partialResults := make([]interface{}, intent.Config.NumWorkers)
	mapIntents := make([]*tasks.Intent, intent.Config.NumWorkers)
	board := generateRandomBoard(1024)
	for i := range partialResults {
		partialResults[i] = GOLPartialResult{
			Data: "Hello",
		}
		mapIntents[i] = &tasks.Intent{
			IntentType: "map",
			TaskType:   intent.TaskType,
			TaskUUID:   intent.TaskUUID,
			Config:     intent.Config,
			Input: map[string]interface{}{
				"size": 1024,
				/*"board": []int8{
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 1, 0, 0, 0,
					0, 0, 0, 0, 0, 1, 0, 0,
					0, 0, 0, 1, 1, 1, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				},*/
				"board": board,
			},
		}
	}

	// Initial state
	state := GOLState{
		round: 0,
	}

	// Create the initial task instance
	taskInstance := tasks.TaskInstance{
		Config:         intent.Config,
		PartialResults: partialResults,
		State:          state,
	}
	return &taskInstance, mapIntents
}

func reduce() *tasks.Intent {
	endIntent := &tasks.Intent{}
	return endIntent
}

func generateRandomBoard(n int) []int8 {
	board := make([]int8, n*n)
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			board[(n*y)+x] = int8(rand.Int() % 2)
		}
	}
	return board
}
