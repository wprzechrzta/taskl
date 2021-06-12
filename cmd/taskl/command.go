package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	task2 "github.com/wprzechrzta/taskl/cmd/taskl/task"
	"log"
)

type ArgRunner interface {
	Init([]string) error
	Run() error
	Name() string
}

type CreateTaskCommand struct {
	fs         *flag.FlagSet
	repository *task2.Repository
	board      string
	body       string
}

func NewTaskCommand(repo *task2.Repository) *CreateTaskCommand {
	tc := &CreateTaskCommand{fs: flag.NewFlagSet("t", flag.PanicOnError), repository: repo}
	tc.fs.StringVar(&tc.board, "b", "My Board", "Board repo attach task")
	return tc
}

func (tc *CreateTaskCommand) Name() string {
	return tc.fs.Name()
}

func (tc *CreateTaskCommand) Init(args []string) error {
	if err := tc.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "Failed repository parse %s", tc.Name())
	}
	if len(tc.fs.Args()) < 1 {
		return fmt.Errorf("TaskComand: Missing task description")
	}
	tc.body = tc.fs.Arg(0)
	return nil
}

func (tc *CreateTaskCommand) Run() error {
	log.Printf("Running command: %s, boardName: %s", tc.Name(), tc.board)
	var t task2.Task
	t.Boards = append(t.Boards, tc.board)
	t.Description = tc.body
	return tc.repository.Create(t)
}
