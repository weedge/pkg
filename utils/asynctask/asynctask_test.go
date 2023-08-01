package asynctask

import (
	"fmt"
	"testing"
)

type testTask struct {
	Name string
}

func newTestTask(Name string) (m *testTask) {
	m = &testTask{
		Name: Name,
	}

	return m
}

func (m *testTask) Run() (err error) {
	//println("name", m.Name)
	if m.Name == "err" {
		return fmt.Errorf("wo qu err le")
	}

	return
}

func TestNewAsyncTask(t *testing.T) {
	asyncTask, err := NewAsyncTask("async.task.test", 3, 3, func(err error) {
	})
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	_, err = NewAsyncTask("async.task.test", 3, 3, func(err error) {
	})
	if err == nil {
		t.Fatal("no error")
		return
	}

	asyncTask.Post(newTestTask("hello"))
	asyncTask.Close()
}

func ExampleAsyncTask_Post() {
	asyncTask, err := NewAsyncTask("async.task.test", 3, 3, func(err error) {
		fmt.Println("Run err", err.Error())
	})
	if err != nil {
		println("err", err.Error())
		return
	}

	asyncTask.Post(newTestTask("hello"), newTestTask("err"))
	asyncTask.Close()

	// Output:
	// Run err wo qu err le
}
