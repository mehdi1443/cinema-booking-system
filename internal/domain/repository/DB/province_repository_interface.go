package DB

import "golang.org/x/net/context"

type ProvinceRepository interface {
	Create(ctx context.Context, province *Province) error
	GetByID(ctx context.Context, id int64) (*Province, error)
	GetByName(ctx context.Context, name string) (*Province, error)
	GetAll(ctx context.Context) ([]Province, error)
}
