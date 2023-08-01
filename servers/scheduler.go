package servers

import (
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
	"time"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	serverId  string

	Tasks           map[string]pufferpanel.Task `json:"tasks"`
	Timezone        string                      `json:"timezone,omitempty"`
	ConcurrentLimit int                         `json:"concurrentLimit"`
	LimitMode       string                      `json:"limitMode"`
}

// LoadScheduler Loads the scheduler from the serverid.cron file, or defaults
// This file is a JSON file, but it hooks into everything
func LoadScheduler(serverId string) (*Scheduler, error) {
	file, err := os.Open(filepath.Join(config.ServersFolder.Value(), serverId+".cron"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	defer pufferpanel.Close(file)

	scheduler := NewDefaultScheduler(serverId)
	if file != nil {
		err = json.NewDecoder(file).Decode(&scheduler)
		if err != nil {
			return nil, err
		}
	}

	err = scheduler.Init()
	return scheduler, err
}

func NewDefaultScheduler(serverId string) *Scheduler {
	return &Scheduler{
		Tasks:           make(map[string]pufferpanel.Task),
		Timezone:        "Local",
		ConcurrentLimit: 5,
		LimitMode:       "wait",
		serverId:        serverId,
	}
}

func (s *Scheduler) Init() error {
	if s.scheduler != nil {
		s.scheduler.Stop()
	}

	var timezone *time.Location
	var err error
	if s.Timezone != "" {
		timezone, err = time.LoadLocation(s.Timezone)
		if err != nil {
			return err
		}
	} else {
		timezone = time.Local
	}

	gs := gocron.NewScheduler(timezone)

	if s.ConcurrentLimit > 0 {
		if s.LimitMode == "reschedule" {
			gs.SetMaxConcurrentJobs(s.ConcurrentLimit, gocron.RescheduleMode)
		} else {
			gs.SetMaxConcurrentJobs(s.ConcurrentLimit, gocron.WaitMode)
		}
	}

	s.scheduler = gs
	return nil
}

func (s *Scheduler) Save() error {
	return nil
}

func (s *Scheduler) Stop() {
	s.scheduler.Stop()
}

func (s *Scheduler) Start() {
	s.scheduler.StartAsync()
}

func (s *Scheduler) IsRunning() bool {
	if s.scheduler == nil {
		return false
	}
	return s.scheduler.IsRunning()
}

func (s *Scheduler) IsTaskRunning(id string) bool {
	jobs, _ := s.scheduler.FindJobsByTag(id)
	for _, v := range jobs {
		if v.IsRunning() {
			return true
		}
	}
	return false
}

func (s *Scheduler) AddTask(id string, task pufferpanel.Task) error {
	job := s.scheduler.Tag(id)
	if task.CronSchedule != "" {
		job = job.Cron(task.CronSchedule)
	}
	_, err := job.Do(_executeTask, s.serverId, id, task)
	return err
}

func (s *Scheduler) RemoveTask(id string) error {
	return s.scheduler.RemoveByTag(id)
}

func (s *Scheduler) RunTask(id string) error {
	return s.scheduler.RunByTag(id)
}

func (s *Scheduler) GetTasks() map[string]pufferpanel.Task {
	return s.Tasks
}

func _executeTask(serverId string, id string, task pufferpanel.Task) {
	p := GetFromCache(serverId)
	var err error

	ops := task.Operations
	if len(ops) > 0 {
		p.RunningEnvironment.DisplayToConsole(true, "Running task %s\n", task.Name)
		var process OperationProcess
		process, err = GenerateProcess(ops, p.GetEnvironment(), p.DataToMap(), p.Execution.EnvironmentVariables)
		if err != nil {
			logging.Error.Printf("Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}

		err = process.Run(p)
		if err != nil {
			logging.Error.Printf("Error setting up tasks: %s", err)
			p.RunningEnvironment.DisplayToConsole(true, "Failed to setup tasks\n")
			p.RunningEnvironment.DisplayToConsole(true, "%s\n", err.Error())
			return
		}
		p.RunningEnvironment.DisplayToConsole(true, "Task %s finished\n", task.Name)
	}
}
