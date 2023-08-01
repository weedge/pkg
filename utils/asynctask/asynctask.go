package asynctask

import (
	"fmt"
	"sync"

	"github.com/weedge/pkg/container/set"
	"github.com/weedge/pkg/safer"
)

var asyncTaskNames *set.HashSet

func init() {
	asyncTaskNames = set.NewSet()
}

type IAsyncTask interface {
	Run() error
}

type AsyncTask struct {
	name string
	ch   chan IAsyncTask
	wg   sync.WaitGroup
}

func (t *AsyncTask) Close() {
	asyncTaskNames.Remove(t.name)
	close(t.ch)
	t.wg.Wait()
}

func GetNamedAsyncTask(name string, taskChanNumber int64, goNumber int, onError func(err error)) (*AsyncTask, error) {
	if ok := asyncTaskNames.Contains(name); ok {
		return nil, fmt.Errorf("asynctask name duplicated: %v", name)
	}
	if taskChanNumber < 1 {
		taskChanNumber = 1
	}
	if goNumber < 1 {
		goNumber = 1
	}
	asyncTaskNames.Add(name)

	asyncTask := new(AsyncTask)
	asyncTask.name = name
	asyncTask.ch = make(chan IAsyncTask, taskChanNumber)
	for i := 0; i < goNumber; i++ {
		safer.GoSafely(
			&asyncTask.wg, false,
			func() { asyncTask.run(onError) },
			nil, nil,
		)
	}
	return asyncTask, nil
}

func (t *AsyncTask) run(onError func(err error)) {
	for {
		asyncTask, ok := <-t.ch
		if !ok {
			break
		}
		if err := asyncTask.Run(); nil != err {
			if nil != onError {
				onError(err)
			}
		}
	}
}

func (t *AsyncTask) Post(tasks ...IAsyncTask) {
	for _, task := range tasks {
		t.ch <- task
	}
}
