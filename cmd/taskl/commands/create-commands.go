package commands

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wprzechrzta/taskl/cmd/taskl/task"
	"log"
)

type CreateTaskCommand struct {
	fs           *flag.FlagSet
	taskOperator *task.TaskOperator
	board        string
	body         string
}

func NewTaskCommand(taskOperator *task.TaskOperator) *CreateTaskCommand {
	tc := &CreateTaskCommand{fs: flag.NewFlagSet("t", flag.PanicOnError), taskOperator: taskOperator}
	tc.fs.StringVar(&tc.board, "b", "My Board", "Board taskOperator attach task")
	return tc
}

func (tc *CreateTaskCommand) Name() string {
	return tc.fs.Name()
}

func (tc *CreateTaskCommand) Init(args []string) error {
	if err := tc.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "Failed taskOperator parse %s", tc.Name())
	}
	if len(tc.fs.Args()) < 1 {
		return fmt.Errorf("TaskComand: Missing task description")
	}
	tc.body = tc.fs.Arg(0)
	return nil
}

func (tc *CreateTaskCommand) Run() error {
	log.Printf("Running command: %s, boardName: %s", tc.Name(), tc.board)
	newTask := createTask(tc.body, tc.board)
	return tc.taskOperator.Create(newTask)

}

func createTask(body string, board string) task.Task {
	var t task.Task
	t.Boards = append(t.Boards, board)
	t.Description = body
	return t
}
