package tmux

func ListSessions(opts ...OptFunc) ([]string, error) {
	return OutputLines("list-sessions", opts...)
}

func WithFormat(format string) OptFunc {
	return WithKeyValue("-F", format)
}

func NewSession(opts ...OptFunc) error {
	return Exec("new-session", opts...)
}

func WithSession(name string) OptFunc {
	return WithKeyValue("-s", name)
}

func WithTarget(target string) OptFunc {
	return WithKeyValue("-t", target)
}

func WithDetached() OptFunc {
	return WithFlag("-d")
}

func SwitchClient(opts ...OptFunc) error {
	return Exec("switch-client", opts...)
}
