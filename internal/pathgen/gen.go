package pathgen

import (
	"crypto/rand"
	"encoding/base64"
)

type Generator struct {
	length int
}

func NewGenerator(length int) *Generator {
	return &Generator{length: length}
}

func (g *Generator) GeneratePath() (string, error) {
	bufLen := (g.length*3 + 3) / 4
	buf := make([]byte, bufLen)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	res := base64.RawURLEncoding.EncodeToString(buf)
	return res[:g.length], nil
}
