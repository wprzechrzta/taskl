package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	task2 "github.com/wprzechrzta/taskl/cmd/taskl/task"
	"log"
	"strings"
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
		return errors.WithMessagef(err, "%s: Failed to fetch tasks ", l.Name())
	}

	sb := strings.Builder{}
	for _, task := range tl.Tasks {
		if _, err = sb.WriteString(fmt.Sprintf("%d. %s \n", task.Id, task.Description)); err != nil {
			return err
		}
	}
	if _, err := sb.WriteString(fmt.Sprintf("Total: %d\n", len(tl.Tasks))); err != nil {
		return err
	}
	fmt.Println(sb.String())
	return err
}

func (l *ListCommand) Name() string {
	return l.fs.Name()
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
