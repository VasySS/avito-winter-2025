package auth

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Repository
type Repository interface {
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}
