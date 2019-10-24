package templates

import (
	"encoding/json"
	"fmt"
	"github.com/pufferpanel/apufferi/v4"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/models"
	"github.com/pufferpanel/pufferpanel/v2/services"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add template from disk",
	Run:   runAdd,
	Args:  cobra.ExactArgs(1),
}

var templateName string
var readme string

func init() {
	AddCmd.Flags().StringVar(&templateName, "name", "", "process id of daemon")
	AddCmd.Flags().StringVar(&readme, "readme", "", "path to readme file")
}

func runAdd(cmd *cobra.Command, args []string) {
	err := pufferpanel.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %s", err.Error())
		os.Exit(1)
		return
	}

	template, err := openTemplate(args[0])

	if err != nil {
		fmt.Printf("Error parsing template: %s", err.Error())
		os.Exit(1)
		return
	}

	if templateName == "" {
		templateName = path.Base(args[0])
	}

	model := &models.Template{
		Template: template,
		Name:     templateName,
		Readme:   "",
	}

	if readme != "" {
		data, err := openReadme(readme)
		if err != nil {
			fmt.Printf("Error reading readme: %s", err.Error())
			os.Exit(1)
			return
		}
		model.Readme = data
	}

	db, err := database.GetConnection()
	if err != nil {
		fmt.Printf("Error getting connection to database: %s", err.Error())
		os.Exit(1)
		return
	}

	ts := &services.Template{DB: db}
	err = ts.Save(model)
	if err != nil {
		fmt.Printf("Error saving template: %s", err.Error())
		os.Exit(1)
		return
	}
	fmt.Printf("Template added: %s", templateName)
}

func openTemplate(path string) (t apufferi.Template, err error) {
	file, err := os.Open(path)
	defer apufferi.Close(file)
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&t)
	return
}

func openReadme(path string) (string, error) {
	file, err := os.Open(path)
	defer apufferi.Close(file)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), err
}
