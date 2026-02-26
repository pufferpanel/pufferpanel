package models

import "golang.org/x/crypto/ssh"

type PublicKey struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement" json:"-"`

	//owners of this permission set
	UserId uint `gorm:"column:user_id;index" json:"-"`
	User   User `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	Name        string `gorm:"column:name;size:100" json:"-" validate:"-"`
	Key         string `gorm:"column:sshkey;size:8000" json:"-" validate:"required"`
	Description string `gorm:"column:description;size:1000" json:"-" validate:"-"`
}

type PublicKeys []*PublicKey

type PublicKeyView struct {
	Key         string `json:"publickey"`
	Name        string `json:"name"`
	Description string `json:"description"`
} //@name PublicKeyView

type PublicKeysView struct {
	Keys []PublicKeyView `json:"keys"`
} //@name PublicKeysView

func (pk *PublicKey) ToView() PublicKeyView {
	return PublicKeyView{
		Key:         pk.Key,
		Name:        pk.Name,
		Description: pk.Description,
	}
}

func (pks PublicKeys) ToView() PublicKeysView {
	keys := make([]PublicKeyView, len(pks))

	for i, v := range pks {
		keys[i] = v.ToView()
	}

	return PublicKeysView{Keys: keys}
}

func (pks PublicKeys) Contains(key string) bool {
	for _, v := range pks {
		if v.Key == key {
			return true
		}
	}

	return false
}

func (pks PublicKeys) ContainsKey(key ssh.PublicKey) bool {
	return pks.Contains(string(key.Marshal()))
}
