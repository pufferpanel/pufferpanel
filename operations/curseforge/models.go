package curseforge

import "time"

type Addon struct {
	Id           uint
	Name         string
	DateCreated  time.Time
	DateModified time.Time
	DateReleased time.Time
	LatestFiles  []File `json:"latestFiles"`
}

type File struct {
	Id                  uint
	IsAvailable         bool
	DisplayName         string
	FileName            string
	ReleaseType         ReleaseType
	FileStatus          int
	FileDate            time.Time
	FileLength          uint64
	DownloadCount       uint
	DownloadUrl         string
	AlternateFileId     uint
	GameVersions        []string
	IsServerPack        bool
	ServerPackFileId    uint
	ParentProjectFileId uint
}

type Category struct {
	Id               uint
	Name             string
	ParentCategoryId uint
	Slug             string
	ClassId          uint
}

type Game struct {
	Id   uint
	Name string
	Slug string
}

type Pagination struct {
	Index       int
	PageSize    int
	ResultCount int
	TotalCount  int
}

func IsAllowedFile(i int) bool {
	switch i {
	case 4:
		fallthrough
	case 10:
		return true
	default:
		return false
	}
}

type Response struct{}

type PagedResponse struct {
	Response
	Pagination Pagination
}

type FilesResponse struct {
	PagedResponse
	Data []File
}

type FileResponse struct {
	Data File
}

type AddonResponse struct {
	Data Addon
}

type Manifest struct {
	Minecraft MinecraftManifest
}

type MinecraftManifest struct {
	Version    string
	ModLoaders []ModLoaderManifest
}

type ModLoaderManifest struct {
	Id      string
	Primary bool
}

type ReleaseType int

var ReleaseFileType = ReleaseType(1)
var BetaFileType = ReleaseType(2)
var AlphaFileType = ReleaseType(3)
