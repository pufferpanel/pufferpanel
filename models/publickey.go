package models

import "golang.org/x/crypto/ssh"

type PublicKey struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement" json:"-"`

	//owners of this permission set
	UserId uint `gorm:"column:user_id;index" json:"-"`
	User   User `gorm:"ASSOCIATION_SAVE_REFERENCE:false" json:"-" validate:"-"`

	Key string `gorm:"column:sshkey;size:8000" json:"-" validate:"-"`
}

type PublicKeys []*PublicKey

type PublicKeyView struct {
	Id  uint   `json:"id"`
	Key string `json:"publickey"`
}

type PublicKeysView struct {
	Keys []string `json:"keys"`
}

func (pk *PublicKey) ToView() PublicKeyView {
	return PublicKeyView{
		Id:  pk.ID,
		Key: pk.Key,
	}
}

func (pks PublicKeys) ToView() PublicKeysView {
	keys := make([]string, len(pks))

	for i, v := range pks {
		keys[i] = v.Key
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
