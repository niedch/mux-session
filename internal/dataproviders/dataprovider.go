package dataproviders

const (
	SELECTED_ICON = ""
	UNSELECTED_ICON = "󰄱"
	WORKTREE_ICON = "󰰱"
)

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
