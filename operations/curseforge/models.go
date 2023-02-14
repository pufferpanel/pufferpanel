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
	Id               uint
	IsAvailable      bool
	DisplayName      string
	FileName         string
	ReleaseType      int
	FileStatus       int
	FileDate         time.Time
	FileLength       uint64
	DownloadCount    uint
	DownloadUrl      string
	AlternateFileId  uint
	GameVersions     []string
	IsServerPack     bool
	ServerPackFileId uint
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

func GetReleaseType(i int) string {
	switch i {
	case 1:
		return "release"
	case 2:
		return "beta"
	case 3:
		return "alpha"
	default:
		return "unknown"
	}
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
