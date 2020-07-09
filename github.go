package pufferpanel

type GithubFileList struct {
	Sha  string       `json:"sha"`
	Url  string       `json:"url"`
	Tree []GithubFile `json:"tree"`
}

type GithubFile struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
	Size uint64 `json:"size"`
	Url  string `json:"url"`
}
