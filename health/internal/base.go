package internal

type usecase struct {
}

func NewUsecase() *usecase {
	return &usecase{}
}

func (u *usecase) CheckHealth() error {
	return nil
}
