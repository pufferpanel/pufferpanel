package programs

import (
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/pufferpanel/pufferpanel/v3"
	"time"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	running   bool
	jobs      map[string]*gocron.Job
	program   *Program
}

// NewScheduler Create a new Scheduler
func NewScheduler(program *Program) Scheduler {
	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(5, gocron.RescheduleMode)
	return Scheduler{scheduler: s, running: false, jobs: map[string]*gocron.Job{}, program: program}
}

// Start the Scheduler (safely)
func (s Scheduler) Start() error {
	if !s.running {
		s.scheduler.StartAsync()
		s.running = true
	}
	return nil
}

func (s Scheduler) isRunning() bool {
	return s.running
}

// Stop the Scheduler if it is running (safely)
func (s Scheduler) Stop() error {
	if s.running {
		s.scheduler.Stop()
		s.running = false
	}
	return nil
}

// Add a single task
func (s Scheduler) Add(task pufferpanel.Task) error {
	var err error
	_, exists := s.jobs[task.Name]
	if exists {
		return errors.New("task already exists")
	}

	job, err := s.scheduler.Cron(task.CronSchedule).Do(s.program.ExecuteTask, task)
	if err != nil {
		return err
	}

	s.jobs[task.Name] = job

	return err
}

// Load a slice of tasks
func (s Scheduler) Load(tasks []pufferpanel.Task) error {
	for _, task := range tasks {
		if err := s.Add(task); err != nil {
			return err
		}
	}
	return nil
}

// LoadMap a map of tasks
func (s Scheduler) LoadMap(tasks map[string]pufferpanel.Task) error {
	for _, task := range tasks {
		if err := s.Add(task); err != nil {
			return err
		}
	}
	return nil
}

// Remove a task from the scheduler
func (s Scheduler) Remove(name string) error {
	j, exists := s.jobs[name]
	if !exists {
		return errors.New("not found")
	}

	s.scheduler.RemoveByReference(j)
	return nil
}

// Rebuild will stop the scheduler, destroy it and create a new instance
func (s Scheduler) Rebuild() error {
	if err := s.Stop(); err != nil {
		return err
	}

	s.scheduler.Clear()
	s.jobs = make(map[string]*gocron.Job)
	return nil
}
