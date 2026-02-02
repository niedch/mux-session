package tmux

func SplitWindow(opts ...OptFunc) error {
	return Exec("split-window", opts...)
}

func WithVertical() OptFunc {
	return WithFlag("-v")
}

func WithHorizontal() OptFunc {
	return WithFlag("-h")
}
