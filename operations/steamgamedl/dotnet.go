//go:build !windows

package steamgamedl

import (
	"os/exec"
)

const DotNetScriptDl = "https://dot.net/v1/dotnet-install.sh"
const DotNetScriptName = "dotnet-install.sh"

func getDotNetInstallCommand() *exec.Cmd {
	//bash dotnet-install.sh -Runtime dotnet -NoPath -InstallDir /var/lib/pufferpanel/binaries/dotnet-runtime
	return exec.Command("bash", "dotnet-install.sh", "-Runtime", "dotnet", "-NoPath", "-InstallDir", "dotnet-runtime")
}
