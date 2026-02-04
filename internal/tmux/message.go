package tmux

func DisplayMessage(opts ...OptFunc) (string, error) {
	return Output("display-message", opts...)
}

func CurrentSession() (string, error) {
	return DisplayMessage(WithPrint(), WithFormat("#S"))
}

func WithPrint() OptFunc {
	return WithFlag("-p")
}

func SendKeys(opts ...OptFunc) error {
	return Exec("send-keys", opts...)
}

func WithKey(key string) OptFunc {
	return WithArg(key)
}

func WithEnter() OptFunc {
	return WithArg("C-m")
}
