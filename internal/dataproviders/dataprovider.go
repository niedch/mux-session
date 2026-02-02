package dataproviders

type Item struct {
	Display string
	Id string
}

type DataProvider interface {
	GetItems() ([]Item, error)
}

