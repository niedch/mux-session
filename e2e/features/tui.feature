Feature: TUI Directory Selection with fzf
  As a user of mux-session
  I want to interactively select directories using fzf
  So that I can quickly create and switch to tmux sessions

  Background:
    Given I build the mux-session

  Scenario: Cancel fzf with Escape key
    Given test directories exist:
      | project-alpha |
      | project-beta  |
    When I spawn mux-session with socket "tui-test-esc" using config:
      """
      search_paths = ["{{TEST_DIR}}"]
      """
    And I press Escape
    Then no tmux session should exist on socket "tui-test-esc"

  Scenario: Cancel fzf with Ctrl+C
    Given test directories exist:
      | project-alpha |
      | project-beta  |
    When I spawn mux-session with socket "tui-test-ctrlc" using config:
      """
      search_paths = ["{{TEST_DIR}}"]
      """
    And I press Ctrl+C
    Then no tmux session should exist on socket "tui-test-ctrlc"

  @manual
  Scenario: Create session by selecting directory with custom window
    This test requires manual interaction as automated PTY input
    to bubbletea/fzf doesn't work reliably in headless environments.
    
    Given test directories exist:
      | project-alpha |
      | project-beta  |
    When I spawn mux-session with socket "tui-test-manual" using config:
      """
      search_paths = ["{{TEST_DIR}}"]

      [[projects]]
      name = "project-alpha"
      dir = "{{TEST_DIR}}/project-alpha"

      [[projects.window_config]]
      window_name = "editor"
      cmd = "echo 'hello'"
      primary = true
      """
    And I manually select "project-alpha" in fzf
    Then a tmux session "project-alpha" should exist on socket "tui-test-manual"
    And session "project-alpha" on socket "tui-test-manual" should have windows:
      | editor |
