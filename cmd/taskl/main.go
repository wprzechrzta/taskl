package main

import (
	"fmt"
	task2 "github.com/wprzechrzta/taskl/cmd/taskl/task"
	"log"
	"os"
)

const defualtStoragePath = "./tmp/.taskl/storage"
const dateFormat = "02 January 2006"

var verbose = false

func Log(fmt string, args ...interface{}) {
	if verbose {
		log.Printf(fmt, args...)
	}
}

func parseAndRun(args []string, config AppConfig) error {
	defaultCommand := "listall"
	if len(args) < 1 {
		args = append(args, defaultCommand)
	}
	taskOperations := task2.NewRepository(config.StoragePath)
	cmds := []ArgRunner{
		NewListCommand(taskOperations),
		NewCreateTaskCommand(taskOperations),
		NewBeginTaskCommand(taskOperations),
		NewCompleteCommand(taskOperations),
		NewCancelCommand(taskOperations),
		NewDeleteCommand(taskOperations),
	}

	subcommand := args[0]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(args[1:]); err != nil {
				return err
			}
			return cmd.Run()
		}
	}
	return fmt.Errorf("Provided subcommand not supported: %s", subcommand)
}

func main() {
	appConfig := AppConfig{StoragePath: defualtStoragePath}
	if err := parseAndRun(os.Args[1:], appConfig); err != nil {
		log.Fatalf("Failed to process request, %+w", err)
	}
}
