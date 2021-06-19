package main

import (
	"fmt"
	"github.com/wprzechrzta/taskl/cmd/taskl/task"
	"io"
	"text/template"
)

type TaskSummary struct {
	Tasks       task.TaskList
	Total       int
	Done        int
	Canceled    int
	Pending     int
	InProgress  int
	DonePercent int
	BoardName   string
}

func calculateSummary(taskList *task.TaskList) (TaskSummary, error) {
	summary := TaskSummary{Tasks: *taskList}
	for _, task := range taskList.Tasks {
		if summary.BoardName == "" {
			summary.BoardName = task.Boards[0]
		}

		if task.IsComplete {
			summary.Done += 1
		} else if task.IsCanelled {
			summary.Canceled += 1
		} else if task.InProgress {
			summary.InProgress += 1
		} else {
			summary.Pending += 1
		}
	}

	summary.Total = len(taskList.Tasks)
	if summary.Total > 0 {
		summary.DonePercent = int((float32(summary.Done) + float32(summary.Canceled)) / float32(summary.Total) * 100)
	}
	return summary, nil
}

func toStatus(task task.Task) string {
	result := "☐"
	if task.InProgress {
		result = "…"
	} else if task.IsCanelled {
		result = "✖"
	} else if task.IsComplete {
		result = "✓"
	}
	return result
}

func renderOutput(out io.Writer, summary TaskSummary) error {
	templ := `{{.BoardName}} [{{ completedTasks .}}/{{.Total}}]
  {{range .Tasks.Tasks}}{{ .Id}}. {{. | toStatus}} {{.Description}} (2 hours)
  {{end}}
{{.Done}} done · {{.Canceled}} canceled · {{.InProgress}} in-progress · {{.Pending}} pending
`

	outputTemplate, err := template.New("output").Funcs(template.FuncMap{
		"toStatus": toStatus,
		"completedTasks": func(summary TaskSummary) int {
			return summary.Done + summary.Canceled
		},
	}).Parse(templ)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return outputTemplate.Execute(out, summary)
}
