package task

import (
	"github.com/google/go-cmp/cmp"
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

func TestShouldRemoveTask(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	expected := TaskList{Tasks: []Task{
		{Id: 1,
			Description: "Example  descriptioin",
			Boards:      []string{"Default Board"},
		},
		{Id: 2,
			Description: "Another task",
			Boards:      []string{"Default Board"},
		},
		{
			Id:          5,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
		},
	}}

	repository := NewRepository(f)

	for _, task := range expected.Tasks {
		_, err = repository.Create(task)
		assert.NoError(t, err)
	}

	saved, err := repository.GetAll()
	assert.Equal(t, len(expected.Tasks), len(saved.Tasks))

	err = repository.Delete(2)

	assert.NoError(t, err)
	saved, err = repository.GetAll()
	assert.Equal(t, 2, len(saved.Tasks))
	//assertContains(t, []Task{expected.Tasks[0], expected.Tasks[2]}, saved.Tasks)

}

func assertContains(t *testing.T, expected, actual []Task) {
	for i := range expected {
		if diff := cmp.Diff(expected[i], actual[i]); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestShouldLoadExistingData(t *testing.T) {
	f, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	expected := TaskList{Tasks: []Task{
		{
			Description: "Example  descriptioin",
			Boards:      []string{"Default Board"},
		},
		{
			Description: "Another task",
			Boards:      []string{"Default Board"},
		},
		{
			Id:          5,
			Description: "Fifth task",
			Boards:      []string{"Default Board"},
		},
	}}

	repository := NewRepository(f)

	loaded, err := repository.GetAll()
	assert.Equal(t, 0, len(loaded.Tasks))

	for _, task := range expected.Tasks {
		_, err = repository.Create(task)
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
	_, err = repository.Create(task)
	assert.NoError(t, err)

	loaded, err = repository.GetAll()
	assert.Equal(t, 1, len(loaded.Tasks))
	assert.Equal(t, 1, loaded.Tasks[0].Id)

	//add task
	task = Task{
		Description: "Second taks",
		Boards:      []string{"Default Board"},
	}
	_, err = repository.Create(task)
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
	_, err = repository.Create(task)

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
	_, err = repository.Create(task)

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
	_, err = repository.Create(task)

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
