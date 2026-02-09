package git

import (
	"os/exec"
	"strings"
)

type Command struct {
	name   string
	args   []string
	output bool
	path   string
}

func (c *Command) Exec() ([]byte, error) {
	args := append([]string{c.name}, c.args...)
	cmd := exec.Command("git", args...)
	if c.path != "" {
		cmd.Dir = c.path
	}
	if c.output {
		return cmd.Output()
	}
	return nil, cmd.Run()
}

type OptFunc func(*Command)

func WithArg(arg string) OptFunc {
	return func(c *Command) {
		c.args = append(c.args, arg)
	}
}

func WithArgs(args ...string) OptFunc {
	return func(c *Command) {
		c.args = append(c.args, args...)
	}
}

func WithFlag(flag string) OptFunc {
	return func(c *Command) {
		c.args = append(c.args, flag)
	}
}

func WithKeyValue(key, value string) OptFunc {
	return func(c *Command) {
		c.args = append(c.args, key, value)
	}
}

func WithOutput() OptFunc {
	return func(c *Command) {
		c.output = true
	}
}

func WithPath(path string) OptFunc {
	return func(c *Command) {
		c.path = path
	}
}

func NewCommand(name string, opts ...OptFunc) *Command {
	cmd := &Command{
		name: name,
		args: make([]string, 0),
	}
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func Exec(name string, opts ...OptFunc) error {
	cmd := NewCommand(name, opts...)
	_, err := cmd.Exec()
	return err
}

func Output(name string, opts ...OptFunc) (string, error) {
	cmd := NewCommand(name, append(opts, WithOutput())...)
	out, err := cmd.Exec()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func OutputLines(name string, opts ...OptFunc) ([]string, error) {
	out, err := Output(name, opts...)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for line := range strings.SplitSeq(out, "\n") {
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}
