package servers

import (
	"encoding/json"
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"os"
	"path/filepath"
	"strings"
)

var (
	allServers = make([]*Server, 0)
)

func LoadFromFolder() {
	err := os.Mkdir(config.ServersFolder.Value(), 0755)
	if err != nil && !os.IsExist(err) {
		logging.Error.Fatalf("Error creating server data folder: %s", err)
	}
	programFiles, err := os.ReadDir(config.ServersFolder.Value())
	if err != nil {
		logging.Error.Fatalf("Error reading from server data folder: %s", err)
	}
	var program *Server
	for _, element := range programFiles {
		if element.IsDir() || !strings.HasSuffix(element.Name(), ".json") {
			continue
		}
		logging.Info.Printf("Attempting to load " + element.Name())
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		program, err = Load(id)
		if err != nil {
			logging.Error.Printf("Error loading server details from json (%s): %s", element.Name(), err)
			continue
		}

		logging.Info.Printf("Loaded server %s", program.Id())
		allServers = append(allServers, program)
	}
}

func GetAll() []*Server {
	return allServers
}

func Load(id string) (program *Server, err error) {
	var data []byte
	data, err = os.ReadFile(filepath.Join(config.ServersFolder.Value(), id+".json"))
	if len(data) == 0 || err != nil {
		return
	}
	program, err = LoadFromData(id, data)
	return
}

func LoadFromData(id string, source []byte) (*Server, error) {
	data := CreateProgram()

	//HACK: Because golang thinks environment and Environment in the json are the same, we have to manually clean the
	//invalid record up....
	rawMap := make(map[string]interface{})
	err := json.Unmarshal(source, &rawMap)
	if err != nil {
		return nil, err
	}

	delete(rawMap, "Environment")
	source, err = json.Marshal(rawMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(source, &data)
	if err != nil {
		return nil, err
	}

	data.Identifier = id

	environmentType := data.Environment.Type

	if environmentType == "standard" || environmentType == "tty" {
		data.Environment.Type = "host"

		err = data.Save()
		if err != nil {
			return nil, err
		}
	}

	data.RunningEnvironment, err = CreateEnvironment(environmentType, config.ServersFolder.Value(), id, data.Environment)
	if err != nil {
		return nil, err
	}

	data.Scheduler, _ = LoadScheduler(data.Id())
	if data.Scheduler == nil {
		data.Scheduler = NewDefaultScheduler(data.Id())
		data.Scheduler.Init()
	}

	return data, nil
}

func Create(program *Server) error {
	if GetFromCache(program.Id()) != nil {
		return pufferpanel.ErrServerAlreadyExists
	}

	var err error

	defer func() {
		if err != nil {
			//revert since we have an error
			_ = os.Remove(filepath.Join(config.ServersFolder.Value(), program.Id()+".json"))
			if program.RunningEnvironment != nil {
				_ = program.RunningEnvironment.Delete()
			}
		}
	}()

	err = program.Save()
	if err != nil {
		logging.Error.Printf("Error writing server: %s", err)
		return err
	}

	replacement, err := Load(program.Id())
	if err != nil {
		return err
	}

	err = replacement.GetEnvironment().Create()
	if err != nil {
		return err
	}

	allServers = append(allServers, replacement)
	return err
}

func Delete(id string) (err error) {
	var index int
	var program *Server
	for i, element := range allServers {
		if element.Id() == id {
			program = element
			index = i
			break
		}
	}
	if program == nil {
		return
	}
	running, err := program.IsRunning()

	if err != nil {
		return
	}

	if running {
		err = program.Kill()
		if err != nil {
			return
		}

		err = program.RunningEnvironment.WaitForMainProcess()
		if err != nil {
			return
		}
	}

	program.Scheduler.Stop()

	err = program.Destroy()
	if err != nil {
		return
	}
	err = os.Remove(filepath.Join(config.ServersFolder.Value(), program.Id()+".json"))
	if err != nil {
		logging.Error.Printf("Error removing server: %s", err)
	}
	allServers = append(allServers[:index], allServers[index+1:]...)
	return
}

func GetFromCache(id string) *Server {
	for _, element := range allServers {
		if element != nil && element.Id() == id {
			return element
		}
	}
	return nil
}

func Save(id string) (err error) {
	program := GetFromCache(id)
	if program == nil {
		err = errors.New("no server with given id")
		return
	}
	err = program.Save()
	return
}

func Reload(id string) (err error) {
	program := GetFromCache(id)
	if program == nil {
		err = errors.New("server does not exist")
		return
	}

	logging.Info.Printf("Reloading server %s", program.Id())
	newVersion, err := Load(id)
	if err != nil {
		logging.Error.Printf("Error reloading server: %s", err)
		return
	}

	program.RunningEnvironment = newVersion.RunningEnvironment
	program.Server = newVersion.Server

	program.Scheduler.Stop()
	logging.Debug.Println("Rebuilding scheduler")

	program.Scheduler, err = LoadScheduler(program.Id())
	if err != nil {
		return err
	}
	if program.Scheduler == nil {
		program.Scheduler = NewDefaultScheduler(program.Id())
		program.Scheduler.Init()
	}

	logging.Debug.Println("Starting scheduler")
	newVersion.Scheduler.Start()
	if err != nil {
		return err
	}

	return
}
