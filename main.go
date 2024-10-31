package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	ID        int
	Job       <-chan string
	Done      chan struct{}
	WaitGroup *sync.WaitGroup
}

type WorkerPool struct {
	Job       chan string
	WaitGroup *sync.WaitGroup
	Workers   []*Worker
}

func NewWorker(id int, job <-chan string, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:        id,
		Job:       job,
		Done:      make(chan struct{}),
		WaitGroup: wg,
	}
}

func (worker *Worker) Work() {
	go func() {
		defer worker.WaitGroup.Done()
		for {
			select {
			case job, ok := <-worker.Job:
				if !ok {
					fmt.Printf("Воркер %d завершил свою работу\n", worker.ID)
					return
				}
				fmt.Printf("Воркер %d выполняет %s\n", worker.ID, job)
				time.Sleep(time.Second)
			case <-worker.Done:
				fmt.Printf("Воркер %d завершил свою работу\n", worker.ID)
				return
			}
		}
	}()
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		Job:       make(chan string, 10),
		WaitGroup: &sync.WaitGroup{},
		Workers:   []*Worker{},
	}
}

func (wp *WorkerPool) InitialWorkers(number int) {
	for i := 0; i < number; i++ {
		wp.WaitGroup.Add(1)
		worker := NewWorker(i, wp.Job, wp.WaitGroup)
		worker.Work()
		wp.Workers = append(wp.Workers, worker)
	}
}

func (wp *WorkerPool) AddWorkers(number int) {
	for i := 0; i < number; i++ {
		wp.WaitGroup.Add(1)
		worker := NewWorker(len(wp.Workers), wp.Job, wp.WaitGroup)
		worker.Work()
		wp.Workers = append(wp.Workers, worker)
	}
}

func (wp *WorkerPool) RemoveWorkers(number int) {
	if number > len(wp.Workers) {
		number = len(wp.Workers)
	}
	for i := 0; i < number; i++ {
		worker := wp.Workers[len(wp.Workers)-1]
		close(worker.Done)
		wp.Workers = wp.Workers[:len(wp.Workers)-1]
	}
}

func main() {
	const initialNumberOfWorkers = 5
	workerPool := NewWorkerPool()
	// Запускаем начальное количество воркеров
	workerPool.InitialWorkers(initialNumberOfWorkers)
	// Отправляем задачи в канал
	go func() {
		for i := 1; i <= 80; i++ {
			workerPool.Job <- fmt.Sprintf("работу номер %d", i)
		}
		close(workerPool.Job)
	}()
	// Динамическое добавление воркеров
	time.Sleep(3 * time.Second) // Ждем немного перед добавлением новых воркеров

	fmt.Println("Добавление 5 новых работников")
	workerPool.AddWorkers(5)

	time.Sleep(3 * time.Second)
	fmt.Println("Удаление 3 работников")
	workerPool.RemoveWorkers(3)
	// Ждем завершения всех воркеров
	workerPool.WaitGroup.Wait()
	fmt.Println("Все работы завершены")
}
