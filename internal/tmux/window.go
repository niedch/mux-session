package tmux

func NewWindow(opts ...OptFunc) error {
	return Exec("new-window", opts...)
}

func SelectWindow(opts ...OptFunc) error {
	return Exec("select-window", opts...)
}

func WithWindowName(name string) OptFunc {
	return WithKeyValue("-n", name)
}

func WithWorkingDir(dir string) OptFunc {
	return WithKeyValue("-c", dir)
}
