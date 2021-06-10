package task

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTaskOperator_Create(t *testing.T) {
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
	}}

	repository := NewRepository(f)

	err = repository.Create(expected.Tasks[0])
	assert.NoError(t, err)

	loaded, err := repository.GetAll()
	assert.Equal(t, len(expected.Tasks), len(loaded.Tasks))

}
