package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	task2 "github.com/wprzechrzta/taskl/cmd/taskl/task"
	"log"
	"os"
	"strconv"
)

type ArgRunner interface {
	Init([]string) error
	Run() error
	Name() string
}

type ListCommand struct {
	fs         *flag.FlagSet
	repository *task2.Repository
}

func NewListCommand(repo *task2.Repository) *ListCommand {
	lc := &ListCommand{fs: flag.NewFlagSet("listall", flag.PanicOnError), repository: repo}
	return lc
}

func (l *ListCommand) Init(args []string) error {
	if err := l.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "%s: Failed parse ", l.Name())
	}
	return nil
}

func (l *ListCommand) Run() error {
	tl, err := l.repository.GetAll()
	if err != nil {
		return errors.WithMessagef(err, "%s: Failed to fetch Tasks ", l.Name())
	}
	summary, err := calculateSummary(tl)
	if err != nil {
		return err
	}
	renderOutput(os.Stdout, summary)
	return err
}

func (l *ListCommand) Name() string {
	return l.fs.Name()
}

//Begin command
type BeginCommand struct {
	fs         *flag.FlagSet
	repository *task2.Repository
	board      string
	taskId     int
}

func NewBeginTaskCommand(repo *task2.Repository) *BeginCommand {
	c := &BeginCommand{fs: flag.NewFlagSet("b", flag.PanicOnError), repository: repo}
	c.fs.StringVar(&c.board, "b", "My Board", "Board repo attach task")
	return c
}

func (b *BeginCommand) Init(args []string) error {
	if err := b.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "%s: Failed parse ", b.Name())
	}

	if len(b.fs.Args()) < 1 {
		return fmt.Errorf("BeginComand: Missing task id")
	}

	taskIdStr := b.fs.Arg(0)
	if taskId, err := strconv.Atoi(taskIdStr); err != nil {
		return errors.WithMessagef(err, "BeginCommand: Task id should be integer value, provided: %v", taskIdStr)
	} else {
		b.taskId = taskId
	}

	return nil
}

func (b *BeginCommand) Run() error {
	return b.repository.Start(b.taskId)
}

func (b *BeginCommand) Name() string {
	return b.fs.Name()
}

//Create command
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
	log.Printf("Running command: %s, BoardName: %s", tc.Name(), tc.board)
	var t task2.Task
	t.Boards = append(t.Boards, tc.board)
	t.Description = tc.body
	return tc.repository.Create(t)
}
