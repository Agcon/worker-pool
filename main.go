package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	ID        int
	Job       <-chan string
	WaitGroup *sync.WaitGroup
}

type WorkerPool struct {
	Job       chan string
	WaitGroup *sync.WaitGroup
}

func NewWorker(id int, job <-chan string, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:        id,
		Job:       job,
		WaitGroup: wg,
	}
}

func (worker *Worker) Work() {
	go func() {
		defer worker.WaitGroup.Done()
		for job := range worker.Job {
			fmt.Printf("Воркер %d выполняет работу %s\n", worker.ID, job)
			time.Sleep(time.Second)
		}
	}()
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		Job:       make(chan string, 10),
		WaitGroup: &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) InitialWorkers(number int) {
	for i := 0; i < number; i++ {
		wp.WaitGroup.Add(1)
		worker := NewWorker(i, wp.Job, wp.WaitGroup)
		worker.Work()
	}
}

func (wp *WorkerPool) AddWorkers(number int) {
	for i := 0; i < number; i++ {
		wp.WaitGroup.Add(1)
		worker := NewWorker(i+1+len(wp.Job), wp.Job, wp.WaitGroup)
		worker.Work()
	}
}

func main() {
	const initialNumberOfWorkers = 5
	workerPool := NewWorkerPool()
	// Запускаем начальное количество воркеров
	workerPool.InitialWorkers(initialNumberOfWorkers)
	// Отправляем задачи в канал
	go func() {
		for i := 1; i <= 50; i++ {
			workerPool.Job <- fmt.Sprintf("Выполнение работы номер %d", i)
		}
		close(workerPool.Job)
	}()
	// Динамическое добавление воркеров
	time.Sleep(3 * time.Second) // Ждем немного перед добавлением новых воркеров

	fmt.Println("Добавление 5 новых работников")
	workerPool.AddWorkers(5)
	// Ждем завершения всех воркеров
	workerPool.WaitGroup.Wait()
	fmt.Println("Все работы завершены")
}
