package main

import (
	"github.com/wprzechrzta/taskl/cmd/taskl/task"
	"io"
)

type TaskSummary struct {
	tasks       task.TaskList
	total       int
	done        int
	canceled    int
	pending     int
	inProgress  int
	donePercent int
	boardName   string
}

func calculateSummary(taskList task.TaskList) (TaskSummary, error) {
	summary := TaskSummary{tasks: taskList}
	for _, task := range taskList.Tasks {
		if summary.boardName == "" {
			summary.boardName = task.Boards[0]
		}

		if task.IsComplete {
			summary.done += 1
		} else if task.IsCanelled {
			summary.canceled += 1
		} else if task.InProgress {
			summary.inProgress += 1
		} else {
			summary.pending += 1
		}
	}

	summary.total = len(taskList.Tasks)
	if summary.total > 0 {
		summary.donePercent = int((float32(summary.done) + float32(summary.canceled)) / float32(summary.total) * 100)
	}
	return summary, nil
}
func renderOutput(out io.Writer, summary TaskSummary) error {

	return nil
}
