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

var verbose = true

func Log(fmt string, args ...interface{}) {
	if verbose {
		log.Printf(fmt, args...)
	}
}

type Task struct {
	Id          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Boards      []string  `json:"boards"`
	InProgress  bool      `json:"inProgress"`
	IsCanelled  bool      `json:"isCanelled"`
	IsComplete  bool      `json:"isComplete"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

type TaskRepository struct {
	StoragePath string
}

func NewRepository(storagePath string) *TaskRepository {
	dbPath := filepath.Join(storagePath, storageFilename)
	err := os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return &TaskRepository{StoragePath: dbPath}
}

func (to *TaskRepository) GetAll() (*TaskList, error) {
	if _, err := os.Stat(to.StoragePath); os.IsNotExist(err) {
		Log("Database file not exists, loc: %v", to.StoragePath)
		return &TaskList{}, nil
	}

	var tasks TaskList
	data, err := ioutil.ReadFile(to.StoragePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to read storage file, loc: %v", to.StoragePath)
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal storage")
	}

	Log("GetAll: %v", tasks)
	return &tasks, nil
}

func (to *TaskRepository) Create(t Task) error {
	Log("Creating task: %+v", t)
	t.Id = nextId()
	t.Date = time.Now()

	//TODO: loadCurrentDB and calculate next id
	allTasks := TaskList{}
	allTasks.Tasks = append(allTasks.Tasks, t)
	data, err := json.MarshalIndent(allTasks, "", " ")
	if err != nil {
		return errors.WithMessage(err, "TaskRepository: Failed to marshal task")
	}
	Log("Create: storing data, loc: %v", to.StoragePath)
	return ioutil.WriteFile(to.StoragePath, data, 0644)
}

func nextId() int {
	return rand.Intn(100)
}
