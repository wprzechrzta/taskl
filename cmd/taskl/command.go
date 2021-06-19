package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	task2 "github.com/wprzechrzta/taskl/cmd/taskl/task"
	"os"
	"strconv"
)

type BasicCommand struct {
	fs         *flag.FlagSet
	repository *task2.Repository
	board      string
	taskId     int
}

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
	*BasicCommand
}

func NewBeginTaskCommand(repo *task2.Repository) *BeginCommand {
	c := &BeginCommand{&BasicCommand{fs: flag.NewFlagSet("b", flag.PanicOnError), repository: repo}}
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
	if err := b.repository.Start(b.taskId); err != nil {
		return err
	}
	fmt.Printf("Started task: %d \n", b.taskId)
	return nil
}

func (b *BeginCommand) Name() string {
	return b.fs.Name()
}

func NewCompleteCommand(repo *task2.Repository) *CompleteCommand {
	c := &CompleteCommand{&BasicCommand{fs: flag.NewFlagSet("c", flag.PanicOnError), repository: repo}}
	c.fs.StringVar(&c.board, "b", "My Board", "Board name tasks belongs to ")
	return c
}

type CompleteCommand struct {
	*BasicCommand
}

func (b *CompleteCommand) Init(args []string) error {
	if err := b.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "%s: Failed parse ", b.Name())
	}

	if len(b.fs.Args()) < 1 {
		return fmt.Errorf("CompleteComand: Missing task id")
	}

	taskIdStr := b.fs.Arg(0)
	if taskId, err := strconv.Atoi(taskIdStr); err != nil {
		return errors.WithMessagef(err, "CompleteCommand: Task id should be integer value, provided: %v", taskIdStr)
	} else {
		b.taskId = taskId
	}
	return nil
}

func (b *CompleteCommand) Run() error {
	if err := b.repository.Complete(b.taskId); err != nil {
		return err
	}
	fmt.Printf("Checked task: %d \n", b.taskId)
	return nil
}

func (b *CompleteCommand) Name() string {
	return b.fs.Name()
}

type CreateTaskCommand struct {
	fs         *flag.FlagSet
	repository *task2.Repository
	board      string
	body       string
}

// NewCreateTaskCommand creates new task
func NewCreateTaskCommand(repo *task2.Repository) *CreateTaskCommand {
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
	Log("Running command: %s, BoardName: %s", tc.Name(), tc.board)
	var t task2.Task
	t.Boards = append(t.Boards, tc.board)
	t.Description = tc.body
	newtask, err := tc.repository.Create(t)
	if err != nil {
		return err
	}
	fmt.Printf("Created task: %d\n", newtask.Id)
	return nil
}

type CancelTaskCommand struct {
	*BasicCommand
}

func NewCancelCommand(repository *task2.Repository) *CancelTaskCommand {
	c := &CancelTaskCommand{&BasicCommand{fs: flag.NewFlagSet("cancel", flag.PanicOnError), repository: repository}}
	c.fs.StringVar(&c.board, "b", "My Board", "Board name tasks belongs to ")
	return c
}

func (c *CancelTaskCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "%s: Failed parse ", c.Name())
	}

	if len(c.fs.Args()) < 1 {
		return fmt.Errorf("CancelComand: Missing task id")
	}

	taskIdStr := c.fs.Arg(0)
	if taskId, err := strconv.Atoi(taskIdStr); err != nil {
		return errors.WithMessagef(err, "CancelCommand: Task id should be integer value, provided: %v", taskIdStr)
	} else {
		c.taskId = taskId
	}
	return nil
}

func (c *CancelTaskCommand) Run() error {
	if err := c.repository.Cancel(c.taskId); err != nil {
		return err
	}
	fmt.Printf("Canceled task: %d \n", c.taskId)
	return nil
}

func (c *CancelTaskCommand) Name() string {
	return c.fs.Name()
}

func NewDeleteCommand(repository *task2.Repository) *DeleteCommand {
	dc := &DeleteCommand{&BasicCommand{fs: flag.NewFlagSet("d", flag.PanicOnError), repository: repository}}
	return dc
}

type DeleteCommand struct {
	*BasicCommand
}

func (d *DeleteCommand) Init(args []string) error {
	if err := d.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "%s: Failed parse ", d.Name())
	}

	if len(d.fs.Args()) < 1 {
		return fmt.Errorf("DeleteComand: Missing task id")
	}

	taskIdStr := d.fs.Arg(0)
	if taskId, err := strconv.Atoi(taskIdStr); err != nil {
		return errors.WithMessagef(err, "DeleteCommand: Task id should be integer value, provided: %v", taskIdStr)
	} else {
		d.taskId = taskId
	}
	return nil
}

func (d *DeleteCommand) Run() error {
	if err := d.repository.Delete(d.taskId); err != nil {
		return err
	}
	fmt.Printf("Deleted task: %d \n", d.taskId)
	return nil
}

func (d *DeleteCommand) Name() string {
	return d.fs.Name()
}
