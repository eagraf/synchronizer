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
	Name:       "GOL",
	Initialize: initialize,
	Setup:      setup,
	Reduce:     reduce,
}

// Setup the game of life task
func setup(intent *tasks.Intent) (*tasks.TaskInstance, []*tasks.Intent) {

	// Initialize partial results and intents
	partialResults := make([]interface{}, intent.Config.NumWorkers)
	mapIntents := make([]*tasks.Intent, intent.Config.NumWorkers)
	board := intent.Input["board"]
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
				"size": intent.Input["size"],
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

// GenerateSetupIntent creates a new setup intent for this type of task
func initialize(UUID string, config tasks.TaskConfig, input map[string]interface{}) *tasks.Intent {
	board := generateRandomBoard(int(input["size"].(float64)))
	newIntent := tasks.Intent{
		IntentType: "setup",
		TaskType:   "GOL",
		TaskUUID:   UUID,
		Config:     config,
		Input: map[string]interface{}{
			"size":  input["size"],
			"board": board,
		},
	}

	return &newIntent
}
