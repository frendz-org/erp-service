package product

type usecase struct {
	repo  ProductRepository
	cache Cache
}

func newUsecase(repo ProductRepository, cache Cache) *usecase {
	return &usecase{
		repo:  repo,
		cache: cache,
	}
}
