//go:build !windows

package steamgamedl

const SteamMetadataLink = SteamMetadataServerLink + "steam_cmd_linux"
const DownloadOs = "linux"

var RenameFolders = map[string]string{
	"linux32": "sdk32",
	"linux64": "sdk64",
}

const DepotDownloaderLink = "https://github.com/SteamRE/DepotDownloader/releases/download/DepotDownloader_2.5.0/DepotDownloader-linux-${arch}.zip"
const DepotDownloaderBinary = "DepotDownloader"
