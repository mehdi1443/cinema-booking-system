package DB

type CityRepository interface {
	Create(ctx context.Context, city *City) error
	GetByID(ctx context.Context, id int64) (*City, error)
	GetByNameAndProvinceID(ctx context.Context, name string, provinceID int64) (*City, error)
	GetByProvinceID(ctx context.Context, provinceID int64) ([]City, error)
	GetAll(ctx context.Context) ([]City, error)
}
