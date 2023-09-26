package messages

type FileList struct {
	CurrentPath string     `json:"path"`
	Error       string     `json:"error,omitempty"`
	Url         string     `json:"url,omitempty"`
	FileList    []FileDesc `json:"files,omitempty"`
	Contents    []byte     `json:"contents,omitempty"`
	Filename    string     `json:"name,omitempty"`
}

type FileDesc struct {
	Name      string `json:"name"`
	Modified  int64  `json:"modifyTime,omitempty"`
	Size      int64  `json:"size,omitempty"`
	File      bool   `json:"isFile"`
	Extension string `json:"extension,omitempty"`
} //@name FileDescription

func (m FileList) Key() string {
	return "file"
}
