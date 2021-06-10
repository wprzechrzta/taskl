package main

import (
	"fmt"
	"github.com/wprzechrzta/taskl/cmd/taskl/commands"
	"github.com/wprzechrzta/taskl/cmd/taskl/env"
	"github.com/wprzechrzta/taskl/cmd/taskl/task"
	"log"
	"os"
)

const defualtStoragePath = "./tmp/.taskl/storage"
const dateFormat = "02 January 2006"

func ParseAndRun(args []string, config env.AppConfig) error {
	if len(args) < 1 {
		log.Println("Running default subcommand: ListCommand")
		return nil
	}
	taskOperations := task.NewRepository(config.StoragePath)

	cmds := []commands.ArgRunner{
		commands.NewTaskCommand(taskOperations),
	}

	subcommand := args[0]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(args[1:])
			return cmd.Run()
		}
	}
	return fmt.Errorf("Provided subcommand not supported: %s", subcommand)
}

func main() {
	appConfig := env.AppConfig{StoragePath: defualtStoragePath}
	log.Println("Starting app...")
	if err := ParseAndRun(os.Args[1:], appConfig); err != nil {
		log.Fatal("Failed to process request, %w", err)
	}
}
