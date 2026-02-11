Feature: Mux Session mark git worktrees
  As a user of mux-session
  I want to see which directories are git worktrees
  So that I can identify worktree directories efficiently

  Scenario: A directory that is a git worktree should be prefixed with "[w]"
    Given a new tmux server
    And I have the following directories:
      | name                |
      | regular-project     |
      | main-repo           |
    And the directory "main-repo" is a git worktree
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    Then I should see the following items in output:
      | item                              |
      | \[ \] .*/regular-project          |
      | \[w\] .*/main-repo                |
      |  └── \[ \] my-worktree            |
      | \[TMUX\] test-session             |
