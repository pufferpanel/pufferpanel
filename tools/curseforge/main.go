package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/operations/curseforge"
	"github.com/pufferpanel/pufferpanel/v3/operations/resolveforgeversion"
	"github.com/pufferpanel/pufferpanel/v3/servers"
)

var tests = []UnitTest{
	{
		CanFail: false,
		Data: curseforge.CurseForge{
			//All the Mods 9 https://www.curseforge.com/minecraft/modpacks/all-the-mods-9/files/7097957
			ProjectId: 715572,
			FileId:    7097957,
		},
	},
	{
		CanFail: false,
		Data: curseforge.CurseForge{
			//Pixelmon https://www.curseforge.com/minecraft/modpacks/the-pixelmon-modpack/files/4966924
			ProjectId: 389615,
			FileId:    4966924,
		},
	},
	{
		CanFail: true,
		Data: curseforge.CurseForge{
			//RLCraft https://www.curseforge.com/minecraft/modpacks/rlcraft/files/4612990
			ProjectId: 285109,
			FileId:    4612990,
		},
	},
	{
		CanFail: false,
		Data: curseforge.CurseForge{
			//Better MC [FABRIC] https://www.curseforge.com/minecraft/modpacks/better-mc-fabric-bmc1/files/4883129
			ProjectId: 452013,
			FileId:    4883129,
		},
	},
	{
		CanFail: false,
		Data: curseforge.CurseForge{
			ProjectId: 876781,
			FileId:    0,
		},
	},
	{
		CanFail: false,
		Data: curseforge.CurseForge{
			//MeatballCraft https://www.curseforge.com/minecraft/modpacks/meatballcraft/files/5842863
			ProjectId: 411966,
			FileId:    5842863,
		},
	},
	{
		CanFail: false,
		Data: curseforge.CurseForge{
			//Farmopolis https://www.curseforge.com/minecraft/modpacks/farmopolis/files/7112573
			ProjectId: 1270262,
			FileId:    7112573,
		},
	},
}

func main() {
	if config.CurseForgeKey.Value() == "" {
		_ = config.CurseForgeKey.Set(os.Getenv("CURSEFORGE_KEY"), false)
	}
	if config.CacheFolder.Value() == "" {
		_ = config.CacheFolder.Set(".", false)
	}

	_ = config.ConsoleForward.Set(true, false)
	_ = config.SecurityDisableUnshare.Set(true, false)

	logging.OriginalStdOut = os.Stdout

	var specificId uint
	flag.UintVar(&specificId, "projectId", 0, "Specific project id")
	flag.Parse()

	results := make(map[UnitTest]error)

	for _, unitTest := range tests {
		test := unitTest.Data

		if specificId != 0 {
			if test.ProjectId != specificId {
				continue
			}
		}

		fmt.Printf("Testing %d\n", test.ProjectId)
		if test.JavaBinary == "" {
			test.JavaBinary = "java"
		}
		serverId := fmt.Sprintf("%d-%d", test.ProjectId, test.FileId)

		_ = os.RemoveAll(serverId)
		_ = os.Mkdir(serverId, 0755)

		server := servers.CreateProgram()
		server.Identifier = serverId

		env, err := servers.CreateEnvironment("host", serverId, "", server)
		if err != nil {
			results[unitTest] = err
			continue
		}
		server.RunningEnvironment = env

		fs, err := files.NewFileServer(serverId, os.Getuid(), os.Getgid())
		if err != nil {
			results[unitTest] = err
			continue
		}

		server.SetFileServer(fs)

		arg := pufferpanel.RunOperatorArgs{
			Environment: env,
			Server:      server,
		}

		result := test.Run(arg)
		if result.Error != nil {
			results[unitTest] = result.Error
			continue
		}
		var fi os.FileInfo
		if fi, err = os.Lstat(filepath.Join(server.GetRootDirectory(), "server.jar")); err == nil && !fi.IsDir() {
			results[unitTest] = nil
		} else {
			op := resolveforgeversion.ResolveForgeVersion{OutputVariable: "result"}
			result = op.Run(arg)
			if result.Error != nil && !os.IsNotExist(err) {
				results[unitTest] = result.Error
				continue
			}
			if result.VariableOverrides == nil || result.VariableOverrides["result"] == "" {
				results[unitTest] = errors.New("failed to resolve to specific version based on unix_args.txt")
			} else {
				results[unitTest] = nil
			}
		}
	}

	failed := false
	for k, v := range results {
		fmt.Printf("Project: %d\n", k.Data.ProjectId)
		if v == nil {
			fmt.Println("  Passes")
		} else {
			fmt.Printf("  Fail: %s\n", v)
			if !k.CanFail {
				failed = true
			}
		}
	}

	if failed {
		os.Exit(1)
	}
}

type UnitTest struct {
	Data    curseforge.CurseForge
	CanFail bool
}
