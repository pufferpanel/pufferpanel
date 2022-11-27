package steamgamedl

import (
	"os/exec"
)

const DotNetScriptDl = "https://dot.net/v1/dotnet-install.ps1"
const DotNetScriptName = "dotnet-install.ps1"

func getDotNetInstallCommand() *exec.Cmd {
	//bash dotnet-install.sh -Runtime dotnet -NoPath -InstallDir /var/lib/pufferpanel/binaries/dotnet-runtime
	return exec.Command("powershell", "-File", "dotnet-install.ps1", "-Runtime", "dotnet", "-NoPath", "-InstallDir", "dotnet-runtime")
}
