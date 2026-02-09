package dataproviders

type Item struct {
	Display    string
	Id         string
	Path       string
	SubItems   []Item
	TreeLevel  int
	IsWorktree bool
}

type DataProvider interface {
	GetItems() ([]Item, error)
}
