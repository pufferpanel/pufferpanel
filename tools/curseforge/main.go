package main

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/operations/curseforge"
	"github.com/pufferpanel/pufferpanel/v3/servers"
	"os"
)

var tests = []curseforge.CurseForge{
	{
		//All the Mods 9 https://www.curseforge.com/minecraft/modpacks/all-the-mods-9/files/5016170
		ProjectId: 715572,
		FileId:    5016170,
	},
	{
		//Pixelmon https://www.curseforge.com/minecraft/modpacks/the-pixelmon-modpack/files/4966924
		ProjectId: 389615,
		FileId:    4966924,
	},
	{
		//RLCraft https://www.curseforge.com/minecraft/modpacks/rlcraft/files/4612990
		ProjectId: 285109,
		FileId:    4612990,
	},
	{
		//Better MC [FABRIC] https://www.curseforge.com/minecraft/modpacks/better-mc-fabric-bmc1/files/4883129
		ProjectId: 452013,
		FileId:    4883129,
	},
}

func main() {
	_ = config.CurseForgeKey.Set(os.Getenv("CURSEFORGE_KEY"), false)
	_ = config.ConsoleForward.Set(true, false)

	logging.OriginalStdOut = os.Stdout

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
			panic(err)
		}
		result := test.Run(env)
		if result.Error != nil {
			panic(result.Error)
		}
	}
}
