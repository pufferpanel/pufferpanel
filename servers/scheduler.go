package servers

import (
	"encoding/json"
	"github.com/go-co-op/gocron/v2"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/utils"
	"os"
	"path/filepath"
	"time"
)

type Scheduler struct {
	scheduler gocron.Scheduler
	serverId  string

	Tasks           map[string]pufferpanel.Task `json:"tasks"`
	Timezone        string                      `json:"timezone,omitempty"`
	ConcurrentLimit uint                        `json:"concurrentLimit"`
	LimitMode       string                      `json:"limitMode"`
}

// LoadScheduler Loads the scheduler from the serverid.cron file, or defaults
// This file is a JSON file, but it hooks into everything
func LoadScheduler(serverId string) (*Scheduler, error) {
	file, err := os.Open(filepath.Join(config.ServersFolder.Value(), serverId+".cron"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	defer utils.Close(file)

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

func (s *Scheduler) Save() error {
	file, err := os.OpenFile(filepath.Join(config.ServersFolder.Value(), s.serverId+".cron"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer utils.Close(file)

	err = json.NewEncoder(file).Encode(s)
	return err
}

func (s *Scheduler) Init() error {
	if s.scheduler != nil {
		err := s.scheduler.StopJobs()
		if err != nil {
			return err
		}
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

	opts := []gocron.SchedulerOption{
		gocron.WithLocation(timezone),
	}

	if s.ConcurrentLimit > 0 {
		if s.LimitMode == "reschedule" {
			opts = append(opts, gocron.WithLimitConcurrentJobs(s.ConcurrentLimit, gocron.LimitModeReschedule))
		} else {
			opts = append(opts, gocron.WithLimitConcurrentJobs(s.ConcurrentLimit, gocron.LimitModeWait))
		}
	}

	gs, err := gocron.NewScheduler(opts...)
	if err != nil {
		return err
	}

	s.scheduler = gs

	for k, v := range s.Tasks {
		err = s.addTask(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) Stop() {
	if s.scheduler.Shutdown() == nil {
		s.scheduler = nil
	}
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) IsRunning() bool {
	return s.scheduler != nil
}

func (s *Scheduler) AddTask(id string, task pufferpanel.Task) error {
	if err := s.addTask(id, task); err != nil {
		return err
	}
	return s.Save()
}

func (s *Scheduler) addTask(id string, task pufferpanel.Task) error {
	var opt gocron.JobDefinition

	if task.CronSchedule != "" {
		opt = gocron.CronJob(task.CronSchedule, true)
	} else {
		opt = gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)))
	}

	_, err := s.scheduler.NewJob(opt, gocron.NewTask(_executeTask, s.serverId, id), gocron.WithName(id))
	if err != nil {
		return err
	}

	s.Tasks[id] = task
	return nil
}

func (s *Scheduler) RemoveTask(id string) error {
	for _, v := range s.scheduler.Jobs() {
		if v.Name() == id {
			_ = s.scheduler.RemoveJob(v.ID())
		}
	}
	delete(s.Tasks, id)
	return s.Save()
}

func (s *Scheduler) RunTask(id string) error {
	jobs := s.scheduler.Jobs()
	for _, v := range jobs {
		if v.Name() == id {
			return v.RunNow()
		}
	}
	return gocron.ErrJobNotFound
}

func (s *Scheduler) GetTasks() map[string]pufferpanel.Task {
	return s.Tasks
}

func _executeTask(serverId string, id string) {
	p := GetFromCache(serverId)
	var err error

	task := p.Scheduler.Tasks[id]

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

func (s *Scheduler) GetExecutor() gocron.Scheduler {
	return s.scheduler
}
