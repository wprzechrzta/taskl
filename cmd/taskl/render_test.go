package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/wprzechrzta/taskl/cmd/taskl/task"
	"testing"
)

func TestCalculateSummary(t *testing.T) {
	assert := assert.New(t)
	tasks := task.TaskList{Tasks: []task.Task{
		task.Task{
			Description: "Example task descriptioin",
			Boards:      []string{"Default Board"},
			InProgress:  true,
		},
		task.Task{
			Description: "Another task",
			Boards:      []string{"Default Board"},
			IsCanelled:  true,
		},
		task.Task{
			Id:          5,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
		},
		task.Task{
			Id:          6,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
			IsComplete:  true,
		},
	}}

	result, error := calculateSummary(tasks)
	assert.NoError(error)
	assert.Equal(result.total, len(tasks.Tasks))
	assert.Equal(result.pending, 1)
	assert.Equal(result.done, 1)
	assert.Equal(result.canceled, 1)
	assert.Equal(result.inProgress, 1)
	assert.Equal(result.donePercent, 50)

}
