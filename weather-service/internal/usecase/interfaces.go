package usecase

type Finder interface {
	Execute(query string) (interface{}, error)
}
