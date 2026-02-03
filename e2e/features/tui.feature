Feature: TUI Directory Selection with fzf
  As a user of mux-session
  I want to interactively select directories using fzf
  So that I can quickly create and switch to tmux sessions

  Background:
    Given I build the mux-session

  @manual
  Scenario: Create session by selecting directory with custom window
    This test requires manual interaction as automated PTY input
    to bubbletea/fzf doesn't work reliably in headless environments.

    Given test directories exist:
      | project-alpha |
      | project-beta  |
    And a new tmux server
    When I spawn mux-session using config:
      """
      search_paths = ["{{TEST_DIR}}"]
      
      [[projects]]
      name = "project-alpha"
      
      [[projects.window_config]]
      window_name = "editor"
      cmd = "echo 'hello'"
      primary = true
      """
    When I select "project-alpha"
    Then I expect following sessions:
      """
      project-alpha
      """
