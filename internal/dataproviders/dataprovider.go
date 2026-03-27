package dataproviders

const (
	SELECTED_ICON   = ""
	UNSELECTED_ICON = "󰄱"
	WORKTREE_ICON   = "󰙅"
	TMUX_ICON       = ""
)

type Item struct {
	Display    string
	Id         string
	Path       string
	SubItems   []Item
	TreeLevel  int
	IsWorktree bool
	ParentId   string
}

type DataProvider interface {
	GetItems() ([]Item, error)
}
