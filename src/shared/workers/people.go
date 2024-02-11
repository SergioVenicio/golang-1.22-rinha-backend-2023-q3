package workers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/config"
	"github.com/sergiovenicio/golang-1.22-rinha-backend-2023-q3/src/domain/person"
)

type Job struct {
	Payload *person.Person
}

type JobQueue chan Job

type PeopleWorker struct {
	Pool  chan chan Job
	Queue JobQueue
	Quit  chan bool
	db    *pgxpool.Pool
	cfg   *config.Config
}

func (w *PeopleWorker) Start() {
	persons := make(chan Job)
	pendingJobs := make(chan []Job)
	go w.bootstrap(persons)
	go w.polling(persons, pendingJobs)
	go w.run(pendingJobs)
}

func (w *PeopleWorker) bootstrap(ch chan Job) {
	for {
		w.Pool <- w.Queue
		select {
		case job := <-w.Queue:
			ch <- job
		case <-w.Quit:
			return
		}
	}
}

func (w *PeopleWorker) polling(persons chan Job, pendingJobs chan []Job) {
	tickInsertRate := time.Duration(time.Duration(w.cfg.Worker.IntsertTime) * time.Second)
	ticker := time.NewTicker(tickInsertRate)
	tickInsert := ticker.C
	batch := make([]Job, 0, w.cfg.Worker.BatchMaxSize)
	for {
		select {
		case data := <-persons:
			batch = append(batch, data)
		case <-tickInsert:
			batchLen := len(batch)
			if batchLen > 0 {
				log.Infof("Tick insert (len=%d)", batchLen)
				pendingJobs <- batch
				batch = make([]Job, 0, w.cfg.Worker.BatchMaxSize)
			}
		}
	}
}

func (w *PeopleWorker) run(pendingJobs chan []Job) {
	columns := []string{"id", "nickname", "name", "birthdate", "stack", "search"}
	identifier := pgx.Identifier{"people"}
	for payload := range pendingJobs {
		_, err := w.db.CopyFrom(
			context.Background(),
			identifier,
			columns,
			pgx.CopyFromSlice(len(payload), w.makeCopyFromSlice(payload)),
		)

		if err != nil {
			log.Errorf("Error on insert batch: %v", err)
		}
	}
}

func (PeopleWorker) makeCopyFromSlice(batch []Job) func(i int) ([]interface{}, error) {
	return func(i int) ([]interface{}, error) {
		return []interface{}{
			batch[i].Payload.ID,
			batch[i].Payload.Nickname,
			batch[i].Payload.Name,
			batch[i].Payload.Birthdate,
			batch[i].Payload.StackStr(),
			batch[i].Payload.SearchStr(),
		}, nil
	}
}

func NewPeopleWorker(
	pool chan chan Job,
	db *pgxpool.Pool,
	cfg *config.Config,
) *PeopleWorker {
	return &PeopleWorker{
		Pool:  pool,
		Queue: make(chan Job, cfg.Worker.Size),
		Quit:  make(chan bool),
		db:    db,
		cfg:   cfg,
	}
}

type Dispatcher struct {
	maxWorkers int
	WorkerPool chan chan Job
	jobQueue   chan Job
	db         *pgxpool.Pool
	cfg        *config.Config
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewPeopleWorker(d.WorkerPool, d.db, d.cfg)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.jobQueue {
		jobChannel := <-d.WorkerPool
		jobChannel <- job
	}
}

func NewDispatcher(
	db *pgxpool.Pool,
	cfg *config.Config,
	jobs JobQueue,
) *Dispatcher {
	maxWorkers := cfg.Worker.MaxWorkers
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		jobQueue:   jobs,
		db:         db,
		cfg:        cfg,
	}
}
