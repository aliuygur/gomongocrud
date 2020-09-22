package gomongocrud

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type UUID4 struct{}

func (_ *UUID4) Generate(ctx context.Context) string {
	return uuid.NewV4().String()
}
