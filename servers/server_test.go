package servers

import (
	"encoding/json"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestProgram_valid(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "invalid json",
			data: "{\n  \"type\": \"\",\n  \"id\": \"902cb39e\",\n  \"run\": {\n    \"autostart\": true\n  }\n}",
			want: false,
		},
		{
			name: "valid json",
			data: "{\n  \"type\": \"minecraft-java\",\n  \"display\": \"Vanilla - Minecraft\",\n  \"id\": \"test\",\n  \"install\": [\n    {\n      \"type\": \"mojangdl\",\n      \"version\": \"latest\",\n      \"target\": \"server.jar\"\n    },\n    {\n      \"type\": \"writefile\",\n      \"text\": \"server-ip=${ip}\\nserver-port=${port}\",\n      \"target\": \"server.properties\"\n    }\n  ],\n  \"run\": {\n    \"stop\": \"stop\",\n    \"pre\": [],\n    \"post\": [],\n    \"command\": \"java -Xmx2048M -Dlog4j2.formatMsgNoLookups=true -jar server.jar nogui\"\n  },\n  \"environment\": {\n    \"type\": \"standard\"\n  },\n  \"data\": {\n    \"ip\": {\n      \"value\": \"0.0.0.0\",\n      \"required\": true,\n      \"desc\": \"What IP to bind the server to\",\n      \"display\": \"IP\",\n      \"internal\": false\n    },\n    \"port\": {\n      \"value\": \"25565\",\n      \"required\": true,\n      \"desc\": \"What port to bind the server to\",\n      \"display\": \"Port\",\n      \"internal\": false,\n      \"type\": \"integer\"\n    }\n  }\n}\n",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Server{
				Server: pufferpanel.Server{},
			}
			err := json.Unmarshal([]byte(tt.data), &p.Server)
			if !assert.NoError(t, err) {
				return
			}

			if got := p.valid(); got != tt.want {
				t.Errorf("valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadFromFolder(t *testing.T) {
	serverData := []byte(`{
  "type": "minecraft-java",
  "data": {
    "eula": {
      "type": "boolean",
      "desc": "Do you (or the server owner) agree to the \u003ca href='https://account.mojang.com/documents/minecraft_eula'\u003eMinecraft EULA?\u003c/a\u003e",
      "display": "EULA Agreement (true/false)",
      "required": true,
      "value": false
    },
    "forgebuild": {
      "type": "",
      "desc": "Version of Forge to install (may be located \u003ca href='http://files.minecraftforge.net/#Downloads'\u003ehere\u003c/a\u003e",
      "display": "Forge Version",
      "value": ""
    },
    "ip": {
      "type": "",
      "desc": "What IP to bind the server to",
      "display": "IP",
      "required": true,
      "value": "0.0.0.0"
    },
    "javaversion": {
      "type": "",
      "desc": "Version of Java to use",
      "display": "Java Version",
      "required": true,
      "value": "17"
    },
    "magmatag": {
      "type": "string",
      "desc": "The alphanumeric commit ID for the build (See https://git.magmafoundation.org/magmafoundation/Magma-1-18-x/-/releases for the commit ID's for each release))",
      "display": "Tag",
      "value": "68e21495",
      "userEdit": true
    },
    "memory": {
      "type": "integer",
      "desc": "How much memory in MB to allocate to the Java Heap",
      "display": "Memory (MB)",
      "required": true,
      "value": 1024
    },
    "modlauncher": {
      "type": "",
      "desc": "What mod or plugin launcher should be installed with Minecraft",
      "display": "Mod/Plugin Launcher",
      "value": "",
      "options": [
        {
          "value": "",
          "display": ""
        },
        {
          "value": "fabric",
          "display": "Fabric"
        },
        {
          "value": "forge",
          "display": "MinecraftForge"
        },
        {
          "value": "magma",
          "display": "Magma"
        },
        {
          "value": "mohist",
          "display": "Mohist"
        },
        {
          "value": "paper",
          "display": "Paper"
        },
        {
          "value": "pufferfish",
          "display": "Pufferfish"
        },
        {
          "value": "purpur",
          "display": "Purpur"
        },
        {
          "value": "quilt",
          "display": "Quilt"
        },
        {
          "value": "spigot",
          "display": "Spigot"
        },
        {
          "value": "ftb",
          "display": "Feed The Beast Modpack"
        },
        {
          "value": "curseforge",
          "display": "CurseForge Modpack"
        }
      ]
    },
    "mohistversion": {
      "type": "string",
      "desc": "Mohist version to install (may be located \u003ca href='https://mohistmc.com/download/' target='_blank'\u003ehere\u003c/a\u003e).",
      "display": "Mohist Version",
      "value": "latest",
      "userEdit": true
    },
    "motd": {
      "type": "",
      "desc": "This is the message that is displayed in the server list of the client, below the name. The MOTD does support \u003ca href='https://minecraft.gamepedia.com/Formatting_codes' target='_blank'\u003ecolor and formatting codes\u003c/a\u003e.",
      "display": "MOTD message of the day",
      "required": true,
      "value": "A Minecraft Server\\n\\u00A79 hosted on PufferPanel"
    },
    "paperbuild": {
      "type": "",
      "desc": "Build of Paper to install (\u003ca href='https://papermc.io/downloads'\u003ePaper version build\u003c/a\u003e). Must be specified as a build number, e.g. 484",
      "display": "build",
      "value": "96"
    },
    "port": {
      "type": "integer",
      "desc": "What port to bind the server to",
      "display": "Port",
      "required": true,
      "value": 25565
    },
    "version": {
      "type": "",
      "desc": "Version of Minecraft you wish to install (not all software may respect this value",
      "display": "Version",
      "required": true,
      "value": "latest"
    }
  },
  "display": "Minecraft: Java Edition",
  "environment": {
    "type": "host"
  },
  "install": [
    {
      "if": "javaversion != \"\"",
      "type": "javadl",
      "version": "${javaversion}"
    },
    {
      "if": "modlauncher == \"\"",
      "target": "server.jar",
      "type": "mojangdl",
      "version": "${version}"
    },
    {
      "if": "modlauncher == \"fabric\"",
      "type": "fabricdl"
    },
    {
      "commands": [
        "java${javaversion} -jar fabric-installer.jar server -mcversion ${version} -downloadMinecraft -noprofile"
      ],
      "if": "modlauncher == \"fabric\"",
      "type": "command"
    },
    {
      "if": "modlauncher == \"fabric\"",
      "source": "fabric-server-launch.jar",
      "target": "server.jar",
      "type": "move"
    },
    {
      "if": "modlauncher == \"forge\"",
      "target": "installer.jar",
      "type": "forgedl",
      "version": "${version}"
    },
    {
      "commands": [
        "java${javaversion} -jar installer.jar --installServer"
      ],
      "if": "modlauncher == \"forge\"",
      "type": "command"
    },
    {
      "files": [
        "https://api.magmafoundation.org/api/v2/${minecraftversion}/latest/${tag}/download"
      ],
      "if": "modlauncher == \"magma\"",
      "type": "download"
    },
    {
      "if": "modlauncher == \"magma\"",
      "source": "Magma-*.jar",
      "target": "server.jar",
      "type": "move"
    },
    {
      "files": [
        "https://mohistmc.com/api/${mc-version}/${mohistversion}/download/"
      ],
      "if": "modlauncher == \"mohist\"",
      "type": "download"
    },
    {
      "if": "modlauncher == \"mohist\"",
      "source": "mohist-*-server.jar",
      "target": "server.jar",
      "type": "move"
    },
    {
      "files": "https://api.papermc.io/v2/projects/paper/versions/${version}/builds/${build}/downloads/paper-${version}-${paperbuild}.jar",
      "if": "modlauncher == \"paper\"",
      "type": "download"
    },
    {
      "if": "modlauncher == \"paper\"",
      "source": "paper-*.jar",
      "target": "paper.jar",
      "type": "move"
    },
    {
      "files": [
        "https://maven.quiltmc.org/repository/release/org/quiltmc/quilt-installer/latest/quilt-installer-latest.jar"
      ],
      "if": "modlauncher == \"quilt\"",
      "type": "download"
    },
    {
      "if": "modlauncher == \"quilt\"",
      "source": "quilt-installer-*.jar",
      "target": "quilt-installer.jar",
      "type": "move"
    },
    {
      "commands": [
        "java${javaversion} -jar quilt-installer.jar install server ${version} --download-server --install-dir=."
      ],
      "if": "modlauncher == \"quilt\"",
      "type": "command"
    },
    {
      "if": "modlauncher == \"quilt\"",
      "source": "quilt-server-launch.jar",
      "target": "server.jar",
      "type": "move"
    },
    {
      "files": "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
      "if": "modlauncher == \"spigot\"",
      "type": "download"
    },
    {
      "commands": [
        "java${javaversion} -jar BuildTools.jar --rev ${version}"
      ],
      "if": "modlauncher == \"spigot\"",
      "type": "command"
    },
    {
      "if": "modlauncher == \"spigot\"",
      "source": "spigot-*.jar",
      "target": "server.jar",
      "type": "move"
    },
    {
      "if": "!file_exists(\"server.properties\")",
      "target": "server.properties",
      "text": "server-ip=${ip}\nserver-port=${port}\nmotd=${motd}\n",
      "type": "writefile"
    },
    {
      "target": "eula.txt",
      "text": "eula=${eula}",
      "type": "writefile"
    }
  ],
  "id": "971f62c5",
  "run": {
    "command": [
      {
        "command": "java${javaversion} -Xmx${memory}M -Dterminal.jline=false -Dterminal.ansi=true -Dlog4j2.formatMsgNoLookups=true @libraries/net/minecraftforge/forge/${version}/win_args.txt nogui",
        "if": "modlauncher == \"forge\" \u0026\u0026 os == \"windows\" \u0026\u0026 file_exists(\"libraries/net/minecraftforge/forge/\" + version+ \"/win_args.txt\")"
      },
      {
        "command": "java${javaversion} -Xmx${memory}M -Dterminal.jline=false -Dterminal.ansi=true -Dlog4j2.formatMsgNoLookups=true @libraries/net/minecraftforge/forge/${version}/unix_args.txt nogui",
        "if": "modlauncher == \"forge\" \u0026\u0026 file_exists(\"libraries/net/minecraftforge/forge/\" + version+ \"/unix_args.txt\")"
      },
      {
        "command": "java${javaversion} -Xmx${memory}M -Dterminal.jline=false -Dterminal.ansi=true -Dlog4j2.formatMsgNoLookups=true -jar server.jar"
      }
    ],
    "stop": "stop",
    "stdin": {
      "type": ""
    }
  },
  "requirements": {},
  "groups": [
    {
      "variables": [
        "eula",
        "memory",
        "ip",
        "port",
        "motd",
        "version",
        "javaversion",
        "modlauncher"
      ],
      "string": "",
      "description": "General settings for all servers",
      "order": 1
    },
    {
      "variables": [
        "forgebuild"
      ],
      "string": "",
      "description": "Settings if using MinecraftForge",
      "order": 2
    },
    {
      "variables": [
        "magmatag"
      ],
      "string": "",
      "description": "Settings specific if using Magma",
      "order": 3
    },
    {
      "variables": [
        "mohistversion"
      ],
      "string": "",
      "description": "Settings specific if using Mohist",
      "order": 4
    },
    {
      "variables": [
        "paperbuild"
      ],
      "string": "",
      "description": "Settings specific if using Paper",
      "order": 5
    }
  ]
}`)

	tmpDir, err := os.MkdirTemp("", "puffer")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(tmpDir)

	config.ServersFolder.Set(tmpDir, false)

	err = os.WriteFile(filepath.Join(tmpDir, "loader.json"), serverData, 0644)
	if !assert.NoError(t, err) {
		return
	}

	err = os.Mkdir(filepath.Join(tmpDir, "loader"), 0755)
	if !assert.NoError(t, err) {
		return
	}

	LoadFromFolder()
	if !assert.NotEmpty(t, allServers) {
		return
	}
}
