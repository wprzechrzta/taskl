package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

type StorageEntry struct {
	Id          int
	Timestamp   string
	Date        string
	Description string
	Boards      []string
	InProgress  bool
	IsCanelled  bool
	IsComplete  bool
}

const StoragePath = "~/.taskl/storage.json"
const dateFormat = "02 January 2006"

type TaskCommand struct {
	fs    *flag.FlagSet
	board string
}

func NewTaskCommand() *TaskCommand {
	tc := &TaskCommand{fs: flag.NewFlagSet("t", flag.PanicOnError)}
	tc.fs.StringVar(&tc.board, "b", "My Board", "Board to attach task")
	return tc
}

func (tc *TaskCommand) Name() string {
	return tc.fs.Name()
}

func (tc *TaskCommand) Init(args []string) error {
	if err := tc.fs.Parse(args); err != nil {
		return errors.WithMessagef(err, "Failed to parse %s", tc.Name())
	}
	return nil
}

func (tc *TaskCommand) Run() error {
	log.Printf("Running command: %s, boardName: %s", tc.Name(), tc.board)
	return nil
}

type ArgRunner interface {
	Init([]string) error
	Run() error
	Name() string
}

func parseAndRun(args []string) error {
	if len(args) < 1 {
		log.Println("Running default subcommand: ListCommand")
		return nil
	}
	cmds := []ArgRunner{
		NewTaskCommand(),
	}

	subcommand := args[0]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(args[1:])
			return cmd.Run()
		}
	}
	return fmt.Errorf("Provided subcommand not supported:", subcommand)
}

func main() {
	log.Println("Starting app...")
	if err := parseAndRun(os.Args[1:]); err != nil {
		log.Fatal("Failed to process request, %w", err)
	}
}
