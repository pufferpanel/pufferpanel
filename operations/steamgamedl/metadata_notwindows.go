//go:build !windows

package steamgamedl

const SteamMetadataLink = SteamMetadataServerLink + "steam_cmd_linux"
const DownloadOs = "linux"
const DepotDownloaderLink = "https://github.com/SteamRE/DepotDownloader/releases/download/DepotDownloader_2.7.3/DepotDownloader-linux-${arch}.zip"
const DepotDownloaderBinary = "DepotDownloader"
