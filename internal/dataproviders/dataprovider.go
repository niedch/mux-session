package dataproviders

type Item struct {
	Display string
	Id      string
	Path    string
}

type DataProvider interface {
	GetItems() ([]Item, error)
}
