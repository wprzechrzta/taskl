package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wprzechrzta/taskl/cmd/taskl/task"
	"testing"
)

func TestCalculateSummary(t *testing.T) {
	assert := assert.New(t)
	tasks := task.TaskList{Tasks: []task.Task{
		{
			Description: "Example task descriptioin",
			Boards:      []string{"Default Board"},
			InProgress:  true,
		},
		{
			Description: "Another task",
			Boards:      []string{"Default Board"},
			IsCanelled:  true,
		},
		{
			Id:          5,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
		},
		{
			Id:          6,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
			IsComplete:  true,
		},
	}}

	result, error := calculateSummary(&tasks)
	assert.NoError(error)
	assert.Equal(result.Total, len(tasks.Tasks))
	assert.Equal(result.Pending, 1)
	assert.Equal(result.Done, 1)
	assert.Equal(result.Canceled, 1)
	assert.Equal(result.InProgress, 1)
	assert.Equal(result.DonePercent, 50)

}

func TestRenderTaskList(t *testing.T) {
	assert := assert.New(t)
	tasks := task.TaskList{[]task.Task{
		{
			Id:          1,
			Description: "First task to render",
			Boards:      []string{"Default Board"},
			InProgress:  true,
		},
		{
			Id:          2,
			Description: "Another task",
			Boards:      []string{"Default Board"},
			IsCanelled:  true,
		},
		{
			Id:          5,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
		},
		{
			Id:          6,
			Description: "Last task to render",
			Boards:      []string{"Default Board"},
			IsComplete:  true,
		},
	}}

	summary, error := calculateSummary(&tasks)
	assert.NoError(error)
	var result bytes.Buffer
	renderOutput(&result, summary)
	fmt.Printf("--[\n%s\n]--", result.String())
	assert.Contains(result.String(), "Default Board")

}
