package servers

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/conditions"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/operations/alterfile"
	"github.com/pufferpanel/pufferpanel/v3/operations/archive"
	"github.com/pufferpanel/pufferpanel/v3/operations/command"
	"github.com/pufferpanel/pufferpanel/v3/operations/console"
	"github.com/pufferpanel/pufferpanel/v3/operations/curseforge"
	"github.com/pufferpanel/pufferpanel/v3/operations/dockerpull"
	"github.com/pufferpanel/pufferpanel/v3/operations/download"
	"github.com/pufferpanel/pufferpanel/v3/operations/extract"
	"github.com/pufferpanel/pufferpanel/v3/operations/fabricdl"
	"github.com/pufferpanel/pufferpanel/v3/operations/forgedl"
	"github.com/pufferpanel/pufferpanel/v3/operations/javadl"
	"github.com/pufferpanel/pufferpanel/v3/operations/mkdir"
	"github.com/pufferpanel/pufferpanel/v3/operations/mojangdl"
	"github.com/pufferpanel/pufferpanel/v3/operations/move"
	"github.com/pufferpanel/pufferpanel/v3/operations/sleep"
	"github.com/pufferpanel/pufferpanel/v3/operations/spongedl"
	"github.com/pufferpanel/pufferpanel/v3/operations/steamgamedl"
	"github.com/pufferpanel/pufferpanel/v3/operations/writefile"
	"github.com/spf13/cast"
)

var commandMapping map[string]pufferpanel.OperationFactory

func init() {
	commandMapping = make(map[string]pufferpanel.OperationFactory)

	loadCoreModules()
}

func GenerateProcess(directions []pufferpanel.ConditionalMetadataType, environment pufferpanel.Environment, dataMapping map[string]interface{}, env map[string]string) (OperationProcess, error) {
	dataMap := make(map[string]interface{})
	for k, v := range dataMapping {
		dataMap[k] = v
	}

	dataMap["rootDir"] = environment.GetRootDirectory()
	operationList := make(OperationProcess, 0)
	for _, mapping := range directions {
		mapCopy := make(map[string]interface{})

		//replace tokens
		for k, v := range mapping.Metadata {
			switch r := v.(type) {
			case string:
				{
					mapCopy[k] = pufferpanel.ReplaceTokens(r, dataMap)
				}
			case []string:
				{
					mapCopy[k] = pufferpanel.ReplaceTokensInArr(r, dataMap)
				}
			case map[string]string:
				{
					mapCopy[k] = pufferpanel.ReplaceTokensInMap(r, dataMap)
				}
			case []interface{}:
				{
					//if we can convert this to a string list, we can work with it
					temp := cast.ToStringSlice(r)
					if len(temp) == len(r) {
						mapCopy[k] = pufferpanel.ReplaceTokensInArr(temp, dataMap)
					} else {
						mapCopy[k] = v
					}
				}
			default:
				mapCopy[k] = v
			}
		}

		envMap := pufferpanel.ReplaceTokensInMap(env, dataMap)

		opCreate := pufferpanel.CreateOperation{
			OperationArgs:        mapCopy,
			EnvironmentVariables: envMap,
			DataMap:              dataMap,
		}

		task := &OperationTask{Type: mapping.Type, Operation: opCreate, Condition: mapping.If}
		operationList = append(operationList, task)
	}
	return operationList, nil
}

type OperationProcess []*OperationTask

type OperationTask struct {
	Operation pufferpanel.CreateOperation
	Condition interface{}
	Type      string
}

func (p *OperationProcess) Run(server *Server) error {
	if len(*p) == 0 {
		return nil
	}

	extraData := map[string]interface{}{
		conditions.VariableSuccess: true,
	}

	var firstError error
	for _, v := range *p {
		shouldRun, err := server.RunCondition(v.Condition, extraData)
		if err != nil {
			return err
		}

		if shouldRun {
			factory := commandMapping[v.Type]
			if factory == nil {
				return pufferpanel.ErrMissingFactory
			}
			op, err := factory.Create(v.Operation)
			if err != nil {
				return pufferpanel.ErrFactoryError(v.Type, err)
			}

			result := op.Run(server.RunningEnvironment)
			if result.Error != nil {
				logging.Error.Printf("Error running command: %s", result.Error.Error())
				if firstError == nil {
					firstError = result.Error
					//TODO: Implement success checking more accurately here
					return result.Error
				}
				extraData[conditions.VariableSuccess] = false
			} else {
				extraData[conditions.VariableSuccess] = true
			}

			if result.VariableOverrides != nil {
				for k, val := range result.VariableOverrides {
					variable := server.Variables[k]
					variable.Value = val
					server.Variables[k] = variable
				}
			}
		}
	}
	return firstError
}

func loadCoreModules() {
	commandFactory := command.Factory
	commandMapping[commandFactory.Key()] = commandFactory

	downloadFactory := download.Factory
	commandMapping[downloadFactory.Key()] = downloadFactory

	mkdirFactory := mkdir.Factory
	commandMapping[mkdirFactory.Key()] = mkdirFactory

	moveFactory := move.Factory
	commandMapping[moveFactory.Key()] = moveFactory

	alterFileFactory := alterfile.Factory
	commandMapping[alterFileFactory.Key()] = alterFileFactory

	writeFileFactory := writefile.Factory
	commandMapping[writeFileFactory.Key()] = writeFileFactory

	mojangFactory := mojangdl.Factory
	commandMapping[mojangFactory.Key()] = mojangFactory

	forgeDlFactory := forgedl.Factory
	commandMapping[forgeDlFactory.Key()] = forgeDlFactory

	fabricDlFactory := fabricdl.Factory
	commandMapping[fabricDlFactory.Key()] = fabricDlFactory

	sleepFactory := sleep.Factory
	commandMapping[sleepFactory.Key()] = sleepFactory

	consoleFactory := console.Factory
	commandMapping[consoleFactory.Key()] = consoleFactory

	archiveFactory := archive.Factory
	commandMapping[archiveFactory.Key()] = archiveFactory

	extractFactory := extract.Factory
	commandMapping[extractFactory.Key()] = extractFactory

	steamgamedlFactory := steamgamedl.Factory
	commandMapping[steamgamedlFactory.Key()] = steamgamedlFactory

	javadlFactory := javadl.Factory
	commandMapping[javadlFactory.Key()] = javadlFactory

	curseforgeFactory := curseforge.Factory
	commandMapping[curseforgeFactory.Key()] = curseforgeFactory

	spongedlFactory := spongedl.Factory
	commandMapping[spongedlFactory.Key()] = spongedlFactory

	dockerpullFactory := dockerpull.Factory
	commandMapping[dockerpullFactory.Key()] = dockerpullFactory

}
