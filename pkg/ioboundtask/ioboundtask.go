package ioboundtask

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

var allTask = struct {
	tasks map[int]task
	mutex sync.RWMutex
}{
	tasks: make(map[int]task),
	mutex: sync.RWMutex{},
}

type task struct {
	ID             int        `json:"id"`
	Work           string     `json:"work"`
	Status         statusType `json:"status"`
	DateCreate     time.Time  `json:"dateCreate"`
	ProcessingTime float64    `json:"processingTime"`
	cancelFunc     context.CancelFunc
}

type statusType string

const (
	statusOk         statusType = "Задача выполнена"
	statusProcessing statusType = "Задача в обработке"
)

var nextID = struct {
	id    int
	mutex sync.Mutex
}{
	id:    0,
	mutex: sync.Mutex{},
}

func CreateTask(work string) {
	newTask := task{Work: work}

	nextID.mutex.Lock()
	newTask.ID = nextID.id
	nextID.id++
	nextID.mutex.Unlock()

	newTask.DateCreate = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	newTask.cancelFunc = cancel
	newTask.Status = statusProcessing

	allTask.mutex.Lock()
	allTask.tasks[newTask.ID] = newTask
	allTask.mutex.Unlock()

	go imitationIOBound(ctx, newTask.ID)
}

func imitationIOBound(ctx context.Context, id int) {
	allTask.mutex.RLock()
	workTask, exists := allTask.tasks[id]
	allTask.mutex.RUnlock()

	if !exists {
		return
	}

	workComplete := make(chan int)

	go func() {
		t := time.NewTimer(
			time.Second * time.Duration(rand.IntN(5*60-3*60)+3*60),
		) // рандомное время от 3 до 5 минут

		select {
		case <-t.C:
			workComplete <- 1
			return
		case <-ctx.Done():
			return
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-workComplete:

		workTask.ProcessingTime = time.Now().Sub(workTask.DateCreate).Seconds()
		workTask.Status = statusOk

		allTask.mutex.Lock()
		allTask.tasks[id] = workTask
		allTask.mutex.Unlock()
		return
	}
}

func DeleteTask(id int) {
	allTask.mutex.Lock()
	deleteTask, exists := allTask.tasks[id]

	if !exists {
		return
	}

	if deleteTask.Status == statusProcessing {
		allTask.tasks[id].cancelFunc()
	}

	delete(allTask.tasks, id)
	allTask.mutex.Unlock()
}

func GetAllTask(rw http.ResponseWriter) error {
	encoder := json.NewEncoder(rw)
	allTask.mutex.RLock()
	err := encoder.Encode(allTask.tasks)
	allTask.mutex.RUnlock()
	if err != nil {
		return err
	}
	return nil
}

func GetTask(rw http.ResponseWriter, id int) error {
	encoder := json.NewEncoder(rw)
	allTask.mutex.RLock()
	err := encoder.Encode(allTask.tasks[id])
	allTask.mutex.RUnlock()
	if err != nil {
		return err
	}
	return nil
}
