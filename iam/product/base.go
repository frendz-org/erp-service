package product

type usecase struct {
	repo  ProductRepository
	cache Cache
}

func NewUsecase(repo ProductRepository, cache Cache) Usecase {
	return &usecase{
		repo:  repo,
		cache: cache,
	}
}
