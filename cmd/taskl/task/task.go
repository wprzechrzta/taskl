package task

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
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
	IsCanelled  bool      `json:"isCancelled"`
	IsComplete  bool      `json:"isComplete"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

type Repository struct {
	StoragePath string
}

func NewRepository(storagePath string) *Repository {
	dbPath := filepath.Join(storagePath, storageFilename)
	err := os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{StoragePath: dbPath}
}

type action func(task *Task)

func (to *Repository) update(id int, updateStrategy action) error {
	tl, err := to.GetAll()
	if err != nil {
		return err
	}
	taskIdx := -1
	for idx := range tl.Tasks {
		if tl.Tasks[idx].Id == id {
			taskIdx = idx
			log.Println("found task:", idx)
			break
		}
	}
	if taskIdx < 0 {
		return fmt.Errorf("Update: Task with id: %d does not exists", id)
	}
	updateStrategy(&tl.Tasks[taskIdx])
	return to.save(tl)
}

func (rep *Repository) Start(id int) error {
	return rep.update(id, func(task *Task) {
		log.Printf("Updating task: %+v", task)
		task.InProgress = true
	})
}

func (rep *Repository) Cancel(id int) error {
	return rep.update(id, func(task *Task) {
		task.IsCanelled = true
	})
}
func (rep *Repository) Complete(id int) error {
	return rep.update(id, func(task *Task) {
		task.IsComplete = true
	})
}

func (to *Repository) GetAll() (*TaskList, error) {
	if _, err := os.Stat(to.StoragePath); os.IsNotExist(err) {
		//Log("Database file not exists, loc: %v", to.StoragePath)
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

	//Log("GetAll: %v", tasks)
	return &tasks, nil
}

func (to *Repository) Create(t Task) error {
	Log("Creating task: %+v", t)
	if t.Id < 1 {
		id, err := to.nextId()
		if err != nil {
			return err
		}
		t.Id = id
	}
	t.Date = time.Now()
	allTasks, err := to.GetAll()
	if err != nil {
		return err
	}
	allTasks.Tasks = append(allTasks.Tasks, t)
	data, err := json.MarshalIndent(allTasks, "", " ")
	if err != nil {
		return errors.WithMessage(err, "Repository: Failed to marshal tasks")
	}
	Log("Create: storing data, loc: %v, data: %v", to.StoragePath, allTasks)
	return ioutil.WriteFile(to.StoragePath, data, 0644)
}
func (to *Repository) save(list *TaskList) error {
	data, err := json.MarshalIndent(list, "", " ")
	if err != nil {
		return errors.WithMessage(err, "Repository: Failed to marshal tasks")
	}
	Log("save: storing  data: %v", list)
	return ioutil.WriteFile(to.StoragePath, data, 0644)
}

func (rep *Repository) nextId() (int, error) {
	tl, err := rep.GetAll()
	if err != nil {
		return -1, err
	}
	var max int
	for _, task := range tl.Tasks {
		if max < task.Id {
			max = task.Id
		}
	}
	return max + 1, nil
}
