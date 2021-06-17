package task

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestTaskOperator_CreateEmptyList(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	repository := NewRepository(f)

	tasks, err := repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, TaskList{}, *tasks)
}

func TestShouldLoadExistingData(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	expected := TaskList{Tasks: []Task{
		Task{
			Description: "Example task descriptioin",
			Boards:      []string{"Default Board"},
		},
		Task{
			Description: "Another task",
			Boards:      []string{"Default Board"},
		},
		Task{
			Id:          5,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
		},
	}}

	repository := NewRepository(f)

	loaded, err := repository.GetAll()
	assert.Equal(t, 0, len(loaded.Tasks))

	for _, task := range expected.Tasks {
		err = repository.Create(task)
		assert.NoError(t, err)
	}
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.Equal(t, len(expected.Tasks), len(loaded.Tasks))

}

func Test_AddTaskAndIncrementId(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	task := Task{
		Description: "Example task descriptioin",
		Boards:      []string{"Default Board"},
	}
	repository := NewRepository(f)

	//when
	loaded, err := repository.GetAll()
	assert.Equal(t, 0, len(loaded.Tasks))

	//add task
	err = repository.Create(task)
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.Equal(t, 1, len(loaded.Tasks))
	assert.Equal(t, 1, loaded.Tasks[0].Id)

	//add task
	task = Task{
		Description: "Second taks",
		Boards:      []string{"Default Board"},
	}
	err = repository.Create(task)
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.Equal(t, 2, len(loaded.Tasks))
	assert.Equal(t, 1, loaded.Tasks[0].Id)
	assert.Equal(t, 2, loaded.Tasks[1].Id)

}

func TestRepository_Start(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	task := Task{
		Id:          1,
		Description: "Wake up",
		Boards:      []string{"MyBoard"},
	}
	repository := NewRepository(f)
	err = repository.Create(task)

	loaded, err := repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loaded.Tasks))
	assert.Equal(t, loaded.Tasks[0].InProgress, false)

	err = repository.Start(task.Id)
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, true, loaded.Tasks[0].InProgress)

}

func TestRepository_Cancel(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	task := Task{
		Id:          1,
		Description: "Wake up",
		Boards:      []string{"MyBoard"},
	}
	repository := NewRepository(f)
	err = repository.Create(task)

	loaded, err := repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loaded.Tasks))
	assert.Equal(t, loaded.Tasks[0].IsCanelled, false)

	err = repository.Cancel(task.Id)
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, true, loaded.Tasks[0].IsCanelled)

}
func TestRepository_Done(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	task := Task{
		Id:          1,
		Description: "Wake up",
		Boards:      []string{"MyBoard"},
	}
	repository := NewRepository(f)
	err = repository.Create(task)

	loaded, err := repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loaded.Tasks))
	assert.Equal(t, loaded.Tasks[0].IsComplete, false)

	err = repository.Complete(task.Id)
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, true, loaded.Tasks[0].IsComplete)

}

func TestConvertion(t *testing.T) {
	in := "123"
	val, err := strconv.Atoi(in)
	assert.NoError(t, err)
	log.Println("---")
	log.Println(val)
	log.Println("---")
}
