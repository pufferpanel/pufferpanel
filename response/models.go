package response

import (
	"github.com/pufferpanel/pufferpanel/v2"
)

type Error struct {
	Error *pufferpanel.Error `json:"error"`
}

type Metadata struct {
	Paging *Paging `json:"paging"`
}

type Paging struct {
	Page    uint `json:"page,omitempty"`
	Size    uint `json:"pageSize,omitempty"`
	MaxSize uint `json:"maxSize,omitempty"`
	Total   uint `json:"total,omitempty"`
}

type Empty struct {
}
