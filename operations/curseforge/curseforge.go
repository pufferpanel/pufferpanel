package curseforge

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/operations/forgedl"
	"github.com/pufferpanel/pufferpanel/v3/operations/neoforgedl"
	"github.com/pufferpanel/pufferpanel/v3/utils"
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
//- neoforgedl
//-   same as forge, but screw downloading
//- quilt
//-   will not deal with

var installerRegex = regexp.MustCompile("(neo)?forge-.*-installer.jar")
var errNoFile = errors.New("status code 404")

func (c CurseForge) Run(args pufferpanel.RunOperatorArgs) pufferpanel.OperationResult {
	env := args.Environment

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

	/*if clientFile.Id != 0 {
		logging.Debug.Printf("Downloading modpack from %s\n", serverFile.DownloadUrl)
		env.DisplayToConsole(true, "Downloading modpack from %s", serverFile.DownloadUrl)
		err = downloadModpack(clientFile)
		if err != nil {
			return pufferpanel.OperationResult{Error: err}
		}
	}*/

	serverZipPath := getCacheFilePath(serverFile)
	logging.Debug.Printf("Extracting modpack from %s\n", serverZipPath)
	env.DisplayToConsole(true, "Extracting modpack from %s", serverZipPath)

	singleRoot, _ := files.DetermineIfSingleRoot(serverZipPath)

	err = files.Extract(nil, serverZipPath, env.GetRootDirectory(), "*", singleRoot, nil)
	if err != nil {
		return pufferpanel.OperationResult{Error: err}
	}

	fs := args.Server.GetFileServer()

	//modpack now downloaded and extracted
	//worse case, this is all we could do...
	//best case, we can get the modpack set up how we need it

	//set 1: resolve the pack to a "modloader"
	var modLoader string
	var data = make(map[string]string)
	var jar string
	var vars map[string]string
	var manifest Manifest
	if jar, err = findInstallerJar(env, fs); err == nil {
		logging.Debug.Printf("Found jar: %s\n", jar)
		if strings.HasPrefix(jar, "neoforgedl") {
			modLoader = "neoforgedl"
		} else {
			modLoader = "forge"
		}
		data["jar"] = jar
	} else if vars, err = readVariableFile(serverFile); err == nil {
		logging.Debug.Printf("Reading variables.txt\n")
		modLoader = strings.ToLower(vars["MODLOADER"])
		data["mcVersion"] = vars["MINECRAFT_VERSION"]
		data["version"] = vars["MODLOADER_VERSION"]
		data["installerVersion"] = vars["FABRIC_INSTALLER_VERSION"]
		logging.Debug.Printf("Resolved: %v\n", data)
	} else if manifest, err = getManifest(clientFile); err == nil {
		logging.Debug.Printf("Using manifest: %v\n", manifest.Minecraft)
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
		logging.Debug.Printf("Resolved: %v\n", data)
	} else {
		//give up
		env.DisplayToConsole(true, "Unknown server type. Could not prepare server for actual execution")
		return pufferpanel.OperationResult{Error: nil}
	}

	//we figured out the loader, now to run their "installer"
	switch modLoader {
	case "fabric":
		{
			err = installFabric(env, fs, data, c.JavaBinary)
			if err != nil {
				return pufferpanel.OperationResult{Error: err}
			}
		}
	case "forge":
		fallthrough
	case "neoforge":
		{
			jarFile := data["jar"]
			if jarFile == "" {
				var downloadUrl string
				var installerJar = "installer.jar"
				version := data["version"]

				if modLoader == "neoforge" {
					downloadUrl = replaceTokens(neoforgedl.InstallerUrl, map[string]string{"version": version})
				} else {
					//because forge has the version in the url, handle it
					mcVersion := data["mcVersion"]
					if !strings.HasPrefix(version, mcVersion) {
						version = mcVersion + "-" + version
					}
					downloadUrl = replaceTokens(forgedl.InstallerUrl, map[string]string{"version": version})
				}

				dl, err := pufferpanel.DownloadViaMaven(downloadUrl, env)
				defer utils.Close(dl)
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
				//copy to server
				err = files.WriteFile(dl, filepath.Join(env.GetRootDirectory(), installerJar))
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
				jarFile = installerJar
			}

			err = installViaJar(args.Server, env, jarFile, c.JavaBinary)
			if err != nil {
				return pufferpanel.OperationResult{Error: err}
			}

			//grab the ServerStarter if there isn't a server.jar, just to help out
			//would prefer Forge's variant, but this will do
			runJarFile := "server.jar"
			if _, err = fs.Stat(runJarFile); os.IsNotExist(err) {
				env.DisplayToConsole(true, "Grabbing ServerStarter")
				var cachePath = filepath.Join(config.CacheFolder.Value(), "github.com", "neoforgedl", "serverstarter", NeoForgeServerStarterVersion, "server.jar")
				if _, err = os.Stat(cachePath); os.IsNotExist(err) {
					env.DisplayToConsole(true, "Downloading "+NeoForgeServerStarter)
					err = pufferpanel.DownloadFileToCache(NeoForgeServerStarter, cachePath)
					if err != nil {
						return pufferpanel.OperationResult{Error: err}
					}
				}
				err = files.CopyFile(cachePath, runJarFile)
				if err != nil {
					return pufferpanel.OperationResult{Error: err}
				}
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

func findInstallerJar(env *pufferpanel.Environment, fs files.FileServer) (string, error) {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		if installerRegex.MatchString(v.Name()) {
			return v.Name(), nil
		}
	}
	return "", os.ErrNotExist
}

func installViaJar(server pufferpanel.DaemonServer, env *pufferpanel.Environment, jarFile string, javaBinary string) error {
	//installer found, we will run this one
	result := make(chan int, 1)
	err := env.Execute(pufferpanel.ExecutionData{
		Command: fmt.Sprintf("%s -jar %s --installServer", javaBinary, jarFile),
		Callback: func(exitCode int) {
			result <- exitCode
			env.DisplayToConsole(true, "Installer exit code: %d", exitCode)
		},
		Variables: server.DataToMap(),
	})
	if err != nil {
		return err
	}
	if <-result != 0 {
		return errors.New("failed to run installer")
	}

	fs := server.GetFileServer()

	//delete installer now
	err = fs.Remove(jarFile)
	if err != nil {
		env.DisplayToConsole(true, "Failed to delete installer")
	}
	err = fs.Remove(jarFile + ".log")
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
		if fi, err = fs.Lstat(f); err == nil && !fi.IsDir() {
			err = fs.Rename(f, "server.jar")
			if err != nil {
				return err
			}
		} else if fi, err = fs.Lstat(f); err == nil && !fi.IsDir() {
			err = fs.Rename(f, "server.jar")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func installFabric(env *pufferpanel.Environment, fs files.FileServer, data map[string]string, javaBinary string) error {
	//this is a mess
	//there's 2 options that exist for fabric
	//there is an "improved" launcher, which is just a jar that we need
	//or we have to pull the installer and run it

	//see if improved is available
	fabricUrl := replaceTokens(ImprovedFabricInstallerUrl, data)
	targetFile := "server.jar"

	env.DisplayToConsole(true, "Downloading %s to %s", fabricUrl, targetFile)
	err := downloadFile(fabricUrl, fs, targetFile)
	if err == nil {
		//this was a good file, we got what we need
		return nil
	} else if !errors.Is(err, errNoFile) {
		//we got a 404, so we can't use the improved version at all
		fabricUrl = replaceTokens(FabricInstallerUrl, data)
		targetFile = "fabric-installer.jar"

		env.DisplayToConsole(true, "Downloading %s to %s", fabricUrl, targetFile)
		err = downloadFile(fabricUrl, fs, targetFile)
		if err != nil {
			return err
		}

		//forge installer found, we will run this one
		result := make(chan int, 1)
		err = env.Execute(pufferpanel.ExecutionData{
			Command: fmt.Sprintf("%s -jar fabric-installer server -mcversion %s -loader %s -downloadMinecraft", javaBinary, data["mcVersion"], data["version"]),
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
		err = fs.Remove("fabric-installer.jar")
		if err != nil {
			env.DisplayToConsole(true, "Failed to delete installer")
		}

		//replace jar with the fabric jar
		_ = fs.Remove("server.jar")
		err = fs.Rename("fabric-server-launch.jar", "server.jar")
		return err
	} else {
		return err
	}
}

func downloadFile(url string, fs files.FileServer, target string) error {
	file, err := fs.OpenFile(target, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer utils.Close(file)
	response, err := pufferpanel.Http().Get(url)
	defer utils.CloseResponse(response)
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
