# README Preview Implementation Plan

## Overview
Implement a markdown preview port that displays the README.md file from the selected project directory.

## Files to Create/Modify

### 1. Create: `internal/fzf/preview_port.go`

```go
package fzf

import (
    "os"
    "path/filepath"
    
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/glamour"
)

type previewPort struct {
    viewport viewport.Model
    renderer glamour.TermRenderer
    content  string
    width    int
    height   int
}

// Methods:
// - newPreviewPort(width, height int) (*previewPort, error)
// - LoadReadme(projectPath string) error
// - SetSize(width, height int)
// - Update(msg tea.Msg) tea.Cmd
// - View() string
```

**Key behaviors:**
- Initialize glamour renderer with `glamour.WithAutoStyle()` and word wrap
- LoadReadme looks for `README.md` in projectPath
- If found: read file, render markdown, update viewport content
- If not found: display "No README.md found in [projectPath]"

### 2. Modify: `internal/fzf/app.go`

**Changes to `model` struct:**
```go
type model struct {
    searchPort  *searchPort
    previewPort *previewPort  // Changed from preview *viewport.Model
    selected    *dataproviders.Item
    width       int
    height      int
}
```

**Changes to `initialModel`:**
- Create previewPort with right side dimensions
- Remove preview viewport creation
- Initial README load for first item (if any)

**Changes to `Init`:**
```go
func (m model) Init() tea.Cmd {
    return tea.Batch(
        m.searchPort.textInput.Focus(),
        func() tea.Msg {
            // Load initial README if items exist
            if item := m.searchPort.GetSelected(); item != nil {
                m.previewPort.LoadReadme(item.Path)
            }
            return nil
        },
    )
}
```

**Changes to `Update`:**
- Handle selection changes: when navigating list, call `previewPort.LoadReadme(selectedItem.Path)`
- Pass messages to both ports: `searchPort.Update(msg)` and `previewPort.Update(msg)`
- Handle resize: update both ports' dimensions

**Changes to `View`:**
- Join `searchPort.View()` and `previewPort.View()` horizontally

### 3. Add Dependency: `go.mod`

Add to imports:
```go
import "github.com/charmbracelet/glamour"
```

Run: `go get github.com/charmbracelet/glamour`

## Implementation Steps

1. **Add glamour dependency**
   - Run go get to add to go.mod

2. **Create preview_port.go**
   - Implement previewPort struct
   - Implement constructor with glamour initialization
   - Implement LoadReadme with file reading and markdown rendering
   - Implement Update and View methods
   - Implement SetSize for responsive layout

3. **Update app.go**
   - Update model struct to use previewPort
   - Update initialModel to create previewPort
   - Update Init for initial README loading
   - Update Update to sync preview with selection changes
   - Update View to render both ports

4. **Test**
   - Build and verify no errors
   - Test with projects that have README.md
   - Test with projects that don't have README.md
   - Verify scrolling works in preview

## Considerations

- **Performance**: README files could be large; viewport handles this well
- **Error handling**: Gracefully handle file read errors
- **Async loading**: Consider loading READMEs asynchronously for large directories
- **Caching**: Could cache rendered READMEs for previously viewed items

## Questions/Decisions

1. Should we show a loading state while reading/rendering?
2. Should we support scrolling within the preview?
3. Any specific styling preferences for "No README found" message?
