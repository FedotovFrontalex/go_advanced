package link

import (
	"math/rand/v2"
	"server/internal/stat"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqIndex"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = RandStringRunes(10)
}

var letterRunes = []rune("abcdefghijklmopqrstuvwxyzABCDFGHIJKLMOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}

	return string(b)
}
