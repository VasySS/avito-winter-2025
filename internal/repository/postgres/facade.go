package postgres

type Facade struct {
	*Repository
}

func NewFacade(repo *Repository) *Facade {
	return &Facade{
		Repository: repo,
	}
}
