package task

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const storageFilename = "taskl.json"

type Task struct {
	Id          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Boards      []string  `json:"boards"`
	InProgress  bool      `json:"inProgress"`
	IsCanelled  bool      `json:"isCanelled"`
	IsComplete  bool      `json:"isComplete"`
}

type TaskCollection struct {
	Tasks []Task `json:"tasks"`
}

type TaskOperator struct {
	storagePath string
}

func (to TaskOperator) Create(t Task) error {
	log.Printf("Creating task: %+v", t)
	t.Id = nextId()
	t.Date = time.Now()

	//TODO: loadCurrentDB and calculate next id
	allTasks := TaskCollection{}
	allTasks.Tasks = append(allTasks.Tasks, t)
	data, err := json.MarshalIndent(allTasks, "", " ")
	if err != nil {
		return errors.WithMessage(err, "TaskOperator: Failed to marsha task")
	}
	return ioutil.WriteFile(to.storagePath, data, 0644)
}

func NewTaskOperator(storagePath string) *TaskOperator {
	storLocation := filepath.Join(storagePath, storageFilename)
	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return &TaskOperator{storagePath: storLocation}
}

func nextId() int {
	return rand.Intn(100)
}
