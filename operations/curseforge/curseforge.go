package curseforge

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CurseForge struct {
	ProjectId  uint
	FileId     uint
	JavaBinary string
}

var forgeInstallerRegex = regexp.MustCompile("forge-.*-installer.jar")
var errNoFile = errors.New("status code 404")

var FabricInstallerUrl = "https://maven.fabricmc.net/net/fabricmc/fabric-installer/${FABRIC_INSTALLER_VERSION}/fabric-installer-${FABRIC_INSTALLER_VERSION}.jar"
var ImprovedFabricInstallerUrl = "https://meta.fabricmc.net/v2/versions/loader/${MINECRAFT_VERSION}/${MODLOADER_VERSION}/${FABRIC_INSTALLER_VERSION}/server/jar"
var ForgeInstallerUrl = "ttps://maven.minecraftforge.net/net/minecraftforge/forge/${MINECRAFT_VERSION}-${MODLOADER_VERSION}/forge-${MINECRAFT_VERSION}-${MODLOADER_VERSION}-installer.jar"

func (c CurseForge) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	client := pufferpanel.Http()

	var file *File
	var err error
	if c.FileId == 0 {
		//we need to get the latest file id to do our calls
		files, err := getLatestFiles(client, c.ProjectId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}

		for _, v := range files {
			if !IsAllowedFile(v.FileStatus) {
				continue
			}
			if file == nil {
				file = &v
				continue
			}
			if file.FileDate.Before(v.FileDate) {
				file = &v
				continue
			}
		}

		if file == nil {
			err = errors.New("no files available on CurseForge")
			return pufferpanel.OperationResult{Error: err}
		}
	} else {
		file, err = getFileById(client, c.ProjectId, c.FileId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	if !file.IsServerPack && file.ServerPackFileId != 0 {
		file, err = getFileById(client, c.ProjectId, file.ServerPackFileId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	downloadZip := filepath.Join(env.GetRootDirectory(), "download.zip")
	err = pufferpanel.DownloadFile(file.DownloadUrl, "download.zip", env)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}
	env.DisplayToConsole(true, "Extracting %s", downloadZip)
	err = pufferpanel.ExtractZipIgnoreSingleDir(downloadZip, env.GetRootDirectory())
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	_ = os.Remove(downloadZip)

	//some mod packs use a variables.txt that contains the relevant data. If it exists, we'll pull that and do what we want
	//others just provide the forge installer directly... so we can find and run it
	//otherwise.... give up
	if fi, err := os.Lstat(filepath.Join(env.GetRootDirectory(), "variables.txt")); err == nil && !fi.IsDir() {
		//read file, we need to pull data from it
		var data map[string]string
		data, err = readVariableFile(filepath.Join(env.GetRootDirectory(), "variables.txt"))
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
		switch strings.ToLower(data["MODLOADER"]) {
		case "fabric":
			{
				err = installFabric(env, client, data, c.JavaBinary)
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
			}
		case "forge":
			{
				forgeUrl := replaceTokens(ForgeInstallerUrl, data)
				forgeInstaller := replaceTokens("forge-${MINECRAFT_VERSION}-${MODLOADER_VERSION}-installer.jar", data)

				err = downloadFile(client, forgeUrl, forgeInstaller)
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}

				err = installForgeViaJar(env, forgeInstaller, c.JavaBinary)
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
			}
		default:
			env.DisplayToConsole(true, "Unknown server type. Could not prepare server for actual execution")
		}
	} else if forgeInstaller, err := findForgeInstallerJar(env); err == nil && forgeInstaller != "" {
		err = installForgeViaJar(env, forgeInstaller, c.JavaBinary)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	} else {
		env.DisplayToConsole(true, "Unknown server type. Could not prepare server for actual execution")
	}

	return pufferpanel.OperationResult{Error: nil}
}

func getLatestFiles(client *http.Client, projectId uint) ([]File, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d", projectId)

	response, err := callCurseForge(client, u)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.CloseResponse(response)

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code from CurseForge: %s", response.Status)
	}

	d, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var addon AddonResponse
	err = json.Unmarshal(d, &addon)
	if err != nil {
		return nil, err
	}
	return addon.Data.LatestFiles, err
}

func getFileById(client *http.Client, projectId, fileId uint) (*File, error) {
	u := fmt.Sprintf("https://api.curseforge.com/v1/mods/%d/files/%d", projectId, fileId)

	response, err := callCurseForge(client, u)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.CloseResponse(response)

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("file id %d not found", fileId)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code from CurseForge: %s", response.Status)
	}

	var res FileResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	return &res.Data, err
}

func callCurseForge(client *http.Client, u string) (*http.Response, error) {
	path, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		Method: "GET",
		URL:    path,
		Header: http.Header{},
	}
	request.Header.Add("x-api-key", config.CurseForgeKey.Value())

	response, err := client.Do(request)
	return response, err
}

func findForgeInstallerJar(env pufferpanel.Environment) (string, error) {
	entries, err := os.ReadDir(env.GetRootDirectory())
	if err != nil {
		return "", err
	}

	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		if forgeInstallerRegex.MatchString(v.Name()) {
			return v.Name(), nil
		}
	}
	return "", os.ErrNotExist
}

func installForgeViaJar(env pufferpanel.Environment, jarFile string, javaBinary string) error {
	//forge installer found, we will run this one
	result := make(chan int, 1)
	err := env.Execute(pufferpanel.ExecutionData{
		Command:   javaBinary,
		Arguments: []string{"-jar", jarFile, "--installServer"},
		Callback: func(exitCode int) {
			result <- exitCode
			env.DisplayToConsole(true, "Installer exit code: %d", exitCode)
		},
	})
	if err != nil {
		return err
	}
	if <-result != 0 {
		return errors.New("failed to run forge installer")
	}

	//delete installer now
	err = os.Remove(filepath.Join(env.GetRootDirectory(), jarFile))
	err = os.Remove(filepath.Join(env.GetRootDirectory(), jarFile+".log"))
	if err != nil {
		env.DisplayToConsole(true, "Failed to delete installer")
	}

	//if this is before 1.16, we have a root jar
	//or if there's a shim
	possibleRenames := []string{
		strings.Replace(jarFile, "-installer", "", 1),      //pre 1.17 forge
		strings.Replace(jarFile, "-installer", "-shim", 1), //forge shim
	}

	var fi os.FileInfo
	for _, f := range possibleRenames {
		if fi, err = os.Lstat(filepath.Join(env.GetRootDirectory(), f)); err == nil && !fi.IsDir() {
			err = os.Rename(filepath.Join(env.GetRootDirectory(), f), filepath.Join(env.GetRootDirectory(), "server.jar"))
			if err != nil {
				return err
			}
		} else if fi, err = os.Lstat(filepath.Join(env.GetRootDirectory(), f)); err == nil && !fi.IsDir() {
			err = os.Rename(filepath.Join(env.GetRootDirectory(), f), filepath.Join(env.GetRootDirectory(), "server.jar"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readVariableFile(path string) (map[string]string, error) {
	data := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer pufferpanel.Close(file)

	scanner := bufio.NewScanner(file)
	var txt string
	for scanner.Scan() {
		txt = scanner.Text()
		parts := strings.SplitN(txt, "=", 2)
		data[parts[0]] = strings.Trim(parts[1], "\"")
	}
	return data, scanner.Err()
}

func installFabric(env pufferpanel.Environment, client *http.Client, data map[string]string, javaBinary string) error {
	//this is a mess
	//there's 2 options that exist for fabric
	//there is an "improved" launcher, which is just a jar that we need
	//or we have to pull the installer and run it

	//see if improved is available
	fabricUrl := replaceTokens(ImprovedFabricInstallerUrl, data)
	targetFile := filepath.Join(env.GetRootDirectory(), "server.jar")

	env.DisplayToConsole(true, "Downloading %s to %s", fabricUrl, targetFile)
	err := downloadFile(client, fabricUrl, targetFile)
	if err == nil {
		//this was a good file, we got what we need
		return nil
	} else if !errors.Is(err, errNoFile) {
		//we got a 404, so we can't use the improved version at all
		fabricUrl = replaceTokens(FabricInstallerUrl, data)
		targetFile = filepath.Join(env.GetRootDirectory(), "fabric-installer.jar")

		env.DisplayToConsole(true, "Downloading %s to %s", fabricUrl, targetFile)
		err = downloadFile(client, fabricUrl, targetFile)
		if err != nil {
			return err
		}

		//forge installer found, we will run this one
		result := make(chan int, 1)
		err = env.Execute(pufferpanel.ExecutionData{
			Command:   javaBinary,
			Arguments: []string{"-jar", "fabric-installer", "server", "-mcversion", data["MINECRAFT_VERSION"], "-loader", data["MODLOADER_VERSION"], "-downloadMinecraft"},
			Callback: func(exitCode int) {
				result <- exitCode
				env.DisplayToConsole(true, "Installer exit code: %d", exitCode)
			},
		})
		if err != nil {
			return err
		}
		if <-result != 0 {
			return errors.New("failed to run fabric installer")
		}

		//delete installer now
		err = os.Remove(filepath.Join(env.GetRootDirectory(), "fabric-installer.jar"))
		if err != nil {
			env.DisplayToConsole(true, "Failed to delete installer")
		}

		//replace jar with the fabric jar
		_ = os.Remove(filepath.Join(env.GetRootDirectory(), "server.jar"))
		err = os.Rename(filepath.Join(env.GetRootDirectory(), "fabric-server-launch.jar"), filepath.Join(env.GetRootDirectory(), "server.jar"))
		return err
	} else {
		return err
	}
}

func downloadFile(client *http.Client, url, target string) error {
	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(file)
	response, err := client.Get(url)
	defer pufferpanel.CloseResponse(response)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		return errNoFile
	}
	_, err = io.Copy(file, response.Body)
	return err
}

func replaceTokens(msg string, data map[string]string) string {
	result := msg
	for k, v := range data {
		result = strings.ReplaceAll(result, "${"+k+"}", v)
	}
	return result
}
