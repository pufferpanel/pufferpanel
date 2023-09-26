package javadl

type File struct {
	Binaries    []Binary `json:"binaries"`
	ReleaseName string   `json:"release_name"`
}

type Binary struct {
	Package Package `json:"package"`
}

type Package struct {
	Link string `json:"link"`
}
