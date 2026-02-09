# mux-session

An interactive tmux session manager that allows you to quickly navigate to project directories and create or switch to tmux sessions.

## Features

- **Interactive Directory Selection**: Uses fzf for fast directory navigation from configured search paths
- **Automatic Session Management**: Creates new tmux sessions or switches to existing ones
- **Project-Specific Configuration**: Define custom window layouts and commands per project
- **Environment Variables**: Set project-specific environment variables
- **Default Window Templates**: Set up default window configurations for all projects

## Installation

### Prerequisites

- Go 1.25.5 or later
- tmux
- fzf

### Build from Source

```bash
# Clone the repository
git clone https://github.com/niedch/mux-session.git
cd mux-session
make install
```

**Add following lines to .tmux.config**

```bash
set-option -g default-shell /bin/zsh
bind-key f run-shell "tmux neww mux-session"
```


## Configuration

Create a configuration file at `$XDG_CONFIG/mux-session/config.toml` (typically `~/.config/mux-session/config.toml`).

### Basic Configuration

```toml
# Directories to search for projects
search_paths = [ "/home/nic/projects", "/home/nic/work" ]

# Default window configuration for all projects
[default]
[[default.window]]
window_name = "Editor"
cmd = "vim ."

[[default.window]]
window_name = "Runner"
cmd = """
lazydocker
"""

[[default.window]]
window_name = "Terminal"
cmd = ""
```

### Project-Specific Configuration

You can override default settings for specific projects:

```toml
# Project-specific configuration
[[project]]
name = "mux-session"

[project.env]
FOO = "bar"
BAZ = "qux"

[[project.window]]
window_name = "nvim"
cmd = "vim ."

[[project]]
name = "dotfiles"

[[project.window]]
window_name = "nvim"
primary = true  # This will be the active window when session starts
cmd = "vim ."

[[project.window]]
window_name = "lazydocker"

[[project.window.panel_config]]
panel_direction = "h"

[[project.window.panel_config]]
panel_direction = "h"
cmd = "lazydocker"
```

### Configuration Options

#### Global Settings
- `search_paths`: Array of directories to search for projects

#### Default Section `[default]`
Defines window templates that apply to all projects unless overridden.

#### Project Section `[[project]]`
- `name`: Project name (must match directory name)
- `env`: A map of environment variables to set for the session.

#### Window Section `[[project.window]]` or `[[default.window]]`
- `window_name`: Name of the tmux window
- `cmd`: Command to run in the window (can be multi-line)
- `primary`: If true, this window will be selected when session starts

#### Panel Configuration `[[project.window.panel_config]]`
- `panel_direction`: Panel direction (`h` for horizontal, `v` for vertical)
- `cmd`: Command to run in this panel

## Usage

### Basic Usage

```bash
# Run with default config
mux-session

# Use custom config file
mux-session -f /path/to/config.toml
```

### Commands

- `mux-session` - Interactive session selection and creation
- `mux-session config-validate` - Validate and display current configuration

### How It Works

1. Launches fzf with directories from your configured search paths
2. Select a directory to work with
3. Checks if a tmux session with that directory name already exists
4. If session exists: switches to it
5. If session doesn't exist: creates new session with configured windows

## Examples

### Quick Start

1. Create your config file:
```bash
mkdir -p ~/.config/mux-session
cp config.toml ~/.config/mux-session/
```

2. Edit the config to add your project directories
3. Run `mux-session` and start managing your sessions!

### Typical Workflow

```bash
# Start mux-session
mux-session

# Select your project from fzf interface
# Automatically creates/switches to tmux session with your configured windows
```

## Development

```bash
# Run in development mode
make dev

# Run tests
make test

# Run e2e tests
make e2e

### Available Make Commands

- `make build` - Build the binary to `bin/mux-session`
- `make run` - Build and run the binary
- `make test` - Run tests
- `make e2e` - Run end-to-end tests
- `make clean` - Clean build artifacts
- `make deps` - Download and tidy dependencies
- `make dev` - Run directly without building binary
- `make install` - Install binary globally with `go install`
- `make all` - Run tests then build
```

## TODOs

- [x] Environment Variables support
- [x] Tmux existing sessions
- [x] Full e2e tests
- [ ] Worktree support
- [x] Display Project info in side panel

