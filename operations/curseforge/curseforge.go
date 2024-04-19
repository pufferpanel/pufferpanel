package curseforge

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"io"
	"net/http"
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

//new plan
//download curseforge to cache if not present (curseforge/projectid/fileid)
//extract from cache to server directory
//determine modloader to use (do this via 3 ways)
//- is there a forge installer present
//- is there a variables.txt present
//- what does the client manifest.json indicate
//install desired mod loader
//- forge
//-   download installer to cache if not present
//-   run installer into cache dir
//-   copy directory to server
//- fabric
//-   check if either installer is in the cache
//-   download improved launcher to the cache
//-   download "old" installer to the cache
//-   run installer in cache
//-   copy directory to server
//- neoforge
//-   not dealing with for now
//- quilt
//-   will not deal with

var forgeInstallerRegex = regexp.MustCompile("forge-.*-installer.jar")
var errNoFile = errors.New("status code 404")

var FabricInstallerUrl = "https://maven.fabricmc.net/net/fabricmc/fabric-installer/${installerVersion}/fabric-installer-${installerVersion}.jar"
var ImprovedFabricInstallerUrl = "https://meta.fabricmc.net/v2/versions/loader/${mcVersion}/${version}/${installerVersion}/server/jar"
var ForgeInstallerUrl = "https://maven.minecraftforge.net/net/minecraftforge/forge/${mcVersion}-${version}/forge-${mcVersion}-${version}-installer.jar"
var ForgeInstallerName = "forge-${mcVersion}-${version}-installer.jar"

func (c CurseForge) Run(env pufferpanel.Environment) pufferpanel.OperationResult {
	var clientFile, serverFile File
	var err error
	if c.FileId == 0 {
		//we need to get the latest file id to do our calls
		files, err := getLatestFiles(c.ProjectId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}

		for _, v := range files {
			if !IsAllowedFile(v.FileStatus) {
				continue
			}
			if v.ReleaseType != ReleaseFileType {
				continue
			}
			if serverFile.FileDate.Before(v.FileDate) {
				serverFile = v
				continue
			}
		}

		if serverFile.Id == 0 {
			err = errors.New("no files available on CurseForge")
			return pufferpanel.OperationResult{Error: err}
		}
	} else {
		serverFile, err = getFileById(c.ProjectId, c.FileId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	if !serverFile.IsServerPack && serverFile.ServerPackFileId != 0 {
		clientFile = serverFile
		serverFile, err = getFileById(c.ProjectId, serverFile.ServerPackFileId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	if clientFile.Id == 0 && serverFile.ParentProjectFileId != 0 {
		clientFile, err = getFileById(c.ProjectId, serverFile.ParentProjectFileId)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	if !serverFile.IsServerPack {
		logging.Debug.Printf("File ID %d is not marked as a server pack, will not install\n", serverFile.Id)
		env.DisplayToConsole(true, "File ID %d is not marked as a server pack, will not install\n", serverFile.Id)
		return pufferpanel.OperationResult{Error: errors.New("not server pack")}
	}

	logging.Debug.Printf("Downloading modpack from %s\n", serverFile.DownloadUrl)
	env.DisplayToConsole(true, "Downloading modpack from %s", serverFile.DownloadUrl)
	err = downloadModpack(serverFile)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	if clientFile.Id != 0 {
		logging.Debug.Printf("Downloading client modpack from %s\n", serverFile.DownloadUrl)
		env.DisplayToConsole(true, "Downloading client modpack from %s", serverFile.DownloadUrl)
		err = downloadModpack(clientFile)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}

	serverZipPath := getCacheFilePath(serverFile)
	logging.Debug.Printf("Extracting modpack from %s\n", serverZipPath)
	env.DisplayToConsole(true, "Extracting modpack from %s", serverZipPath)
	err = pufferpanel.Extract(nil, serverZipPath, env.GetRootDirectory(), "*", true)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	//modpack now downloaded and extracted
	//worse case, this is all we could do...
	//best case, we can get the modpack set up how we need it

	//set 1: resolve the pack to a "modloader"
	var modLoader string
	var data = make(map[string]string)
	var jar string
	var vars map[string]string
	var manifest Manifest
	if jar, err = findForgeInstallerJar(env); err == nil {
		modLoader = "forge"
		data["jar"] = jar
	} else if vars, err = readVariableFile(serverFile); err == nil {
		modLoader = strings.ToLower(vars["MODLOADER"])
		data["mcVersion"] = vars["MINECRAFT_VERSION"]
		data["version"] = vars["MODLOADER_VERSION"]
		data["installerVersion"] = vars["FABRIC_INSTALLER_VERSION"]
	} else if manifest, err = getManifest(clientFile); err == nil {
		mcVersion := manifest.Minecraft.Version
		var loaderVersion string
		for _, v := range manifest.Minecraft.ModLoaders {
			if v.Primary {
				loaderVersion = v.Id
				break
			}
		}
		parts := strings.SplitN(loaderVersion, "-", 2)
		modLoader = parts[0]
		data["mcVersion"] = mcVersion
		data["version"] = parts[1]
	} else {
		//give up
		env.DisplayToConsole(true, "Unknown server type. Could not prepare server for actual execution")
		return pufferpanel.OperationResult{Error: nil}
	}

	//we figured out the loader, now to run their "installer"
	switch modLoader {
	case "forge":
		{
			//for forge, we need a jar
			//so, if we don't have one already, we'll have to get it
			jarFile := data["jar"]
			if jarFile == "" {
				forgeUrl := replaceTokens(ForgeInstallerUrl, data)
				forgeInstaller := replaceTokens(ForgeInstallerName, data)

				jarFile, err = pufferpanel.DownloadViaMaven(forgeUrl, env)
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
				//copy to server
				err = pufferpanel.CopyFile(jarFile, filepath.Join(env.GetRootDirectory(), forgeInstaller))
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
				jarFile = forgeInstaller
			}
			err = installForgeViaJar(env, jarFile, c.JavaBinary)
			if err != nil {
				return pufferpanel.OperationResult{Error: err}
			}
		}
	case "fabric":
		{
			err = installFabric(env, data, c.JavaBinary)
			if err != nil {
				return pufferpanel.OperationResult{Error: err}
			}
		}
	default:
		{
			env.DisplayToConsole(true, "Unsupported server type. Could not prepare server for actual execution")
			return pufferpanel.OperationResult{Error: nil}
		}
	}

	//loaders installed, at this stage, we're "done"
	env.DisplayToConsole(true, "Pack installed and should be good to go!")
	return pufferpanel.OperationResult{Error: nil}
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

func installFabric(env pufferpanel.Environment, data map[string]string, javaBinary string) error {
	//this is a mess
	//there's 2 options that exist for fabric
	//there is an "improved" launcher, which is just a jar that we need
	//or we have to pull the installer and run it

	//see if improved is available
	fabricUrl := replaceTokens(ImprovedFabricInstallerUrl, data)
	targetFile := filepath.Join(env.GetRootDirectory(), "server.jar")

	env.DisplayToConsole(true, "Downloading %s to %s", fabricUrl, targetFile)
	err := downloadFile(fabricUrl, targetFile)
	if err == nil {
		//this was a good file, we got what we need
		return nil
	} else if !errors.Is(err, errNoFile) {
		//we got a 404, so we can't use the improved version at all
		fabricUrl = replaceTokens(FabricInstallerUrl, data)
		targetFile = filepath.Join(env.GetRootDirectory(), "fabric-installer.jar")

		env.DisplayToConsole(true, "Downloading %s to %s", fabricUrl, targetFile)
		err = downloadFile(fabricUrl, targetFile)
		if err != nil {
			return err
		}

		//forge installer found, we will run this one
		result := make(chan int, 1)
		err = env.Execute(pufferpanel.ExecutionData{
			Command:   javaBinary,
			Arguments: []string{"-jar", "fabric-installer", "server", "-mcversion", data["mcVersion"], "-loader", data["version"], "-downloadMinecraft"},
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

func downloadFile(url, target string) error {
	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer pufferpanel.Close(file)
	response, err := pufferpanel.Http().Get(url)
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
