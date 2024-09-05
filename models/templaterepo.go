package models

type TemplateRepo struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;not null;size:100" json:"name"`
	Url      string `gorm:"column:url;not null;size:4000" json:"url"`
	Branch   string `gorm:"column:branch;not null;size:100;default:'main'" json:"branch"`
	PAT      string `gorm:"column:pat;size:4000" json:"-"`
	Username string `gorm:"column:username;size:1000" json:"-"`
	Password string `gorm:"column:password;size:1000" json:"-"`
	SSHKey   string `gorm:"column:ssh_key;size:4000" json:"-"`
	IsLocal  bool   `gorm:"-" json:"isLocal"`
} //@name TemplateRepo
