package tmux

func SetEnvironment(opts ...OptFunc) error {
	return Exec("set-environment", opts...)
}
