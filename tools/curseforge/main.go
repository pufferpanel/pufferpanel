package main

import (
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/operations/curseforge"
	"github.com/pufferpanel/pufferpanel/v3/operations/resolveforgeversion"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"os"
	"path/filepath"
)

var tests = []curseforge.CurseForge{
	//{
	//	//All the Mods 9 https://www.curseforge.com/minecraft/modpacks/all-the-mods-9/files/5016170
	//	ProjectId: 715572,
	//	FileId:    5016170,
	//},
	//{
	//	//Pixelmon https://www.curseforge.com/minecraft/modpacks/the-pixelmon-modpack/files/4966924
	//	ProjectId: 389615,
	//	FileId:    4966924,
	//},
	{
		//RLCraft https://www.curseforge.com/minecraft/modpacks/rlcraft/files/4612990
		ProjectId: 285109,
		FileId:    4612990,
	},
	//{
	//	//Better MC [FABRIC] https://www.curseforge.com/minecraft/modpacks/better-mc-fabric-bmc1/files/4883129
	//	ProjectId: 452013,
	//	FileId:    4883129,
	//},
}

func main() {
	_ = config.CurseForgeKey.Set(os.Getenv("CURSEFORGE_KEY"), false)
	_ = config.ConsoleForward.Set(true, false)

	logging.OriginalStdOut = os.Stdout

	results := make(map[uint]error)

	for _, test := range tests {
		fmt.Printf("Testing %d\n", test.ProjectId)
		if test.JavaBinary == "" {
			test.JavaBinary = "java"
		}
		serverId := fmt.Sprintf("%d-%d", test.ProjectId, test.FileId)

		_ = os.RemoveAll(serverId)
		_ = os.Mkdir(serverId, 0755)

		env, err := servers.CreateEnvironment("host", ".", serverId, pufferpanel.MetadataType{Type: "host"})
		if err != nil {
			results[test.ProjectId] = err
			continue
		}
		result := test.Run(env)
		if result.Error != nil {
			results[test.ProjectId] = err
			continue
		}
		var fi os.FileInfo
		if fi, err = os.Lstat(filepath.Join(serverId, "server.jar")); err == nil && !fi.IsDir() {
			results[test.ProjectId] = nil
		} else {
			op := resolveforgeversion.ResolveForgeVersion{OutputVariable: "result"}
			result = op.Run(env)
			if result.Error != nil {
				results[test.ProjectId] = result.Error
				continue
			}
			if result.VariableOverrides["result"] == "" {
				results[test.ProjectId] = errors.New("failed to resolve to specific MC Forge version based on unix_args.txt")
			} else {
				results[test.ProjectId] = nil
			}
		}
	}

	for k, v := range results {
		fmt.Printf("Project: %d\n", k)
		if v == nil {
			fmt.Println("  Passes")
		} else {
			fmt.Printf("  Fail: %s\n", v)
		}
	}
}
