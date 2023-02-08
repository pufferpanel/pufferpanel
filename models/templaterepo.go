package models

type TemplateRepo struct {
	Name     string `gorm:"type:varchar(100);primaryKey" json:"name"`
	Url      string `gorm:"type:text" json:"url"`
	Branch   string `gorm:"type:text" json:"branch"`
	PAT      string `json:"-"`
	Username string `json:"-"`
	Password string `json:"-"`
	SSHKey   string `json:"-"`
}
