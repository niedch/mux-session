Feature: Mux Session mark git worktrees
  As a user of mux-session
  I want to see which directories are git worktrees
  So that I can identify worktree directories efficiently

  Scenario: A directory that is a git worktree should be prefixed with "[w]"
    Given a new tmux server
    And I have the following directories:
      | name            |
      | regular-project |
      | main-repo       |
    And the directory "main-repo" is a git worktree
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    Then I should see the following items in output:
      | item                 |
      |  test-session       |
      | 󰙅 .*/main-repo       |
      | └── 󰄱 my-worktree    |
      | 󰄱 .*/regular-project |

  Scenario: A worktree that has an active session should be marked with both "[w]" and "[x]"
    Given a new tmux server
    And I have the following directories:
      | name            |
      | regular-project |
      | main-repo       |
    And the directory "main-repo" is a git worktree
    When I run mux-session switch "my-worktree" with config:
      """
      search_paths = ["<search_path>"]
      
      [default]
      [[default.window]]
      window_name = "Shell"
      """
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    Then I should see the following items in output:
      | item                 |
      |  test-session       |
      | 󰙅 .*/main-repo       |
      | └── .*my-worktree     |
      | 󰄱 .*/regular-project |

  Scenario: A worktree should inherit its parent's window configuration
    Given a new tmux server
    And I have the following directories:
      | name            |
      | main-repo       |
    And the directory "main-repo" is a git worktree
    When I run mux-session switch "my-worktree" with config:
      """
      search_paths = ["<search_path>"]
      
      [default]
      [[default.window]]
      window_name = "DefaultShell"

      [[project]]
      name = "main-repo"
      
      [[project.window]]
      window_name = "WorktreeInherited"
      """
    Then I expect following sessions:
      """
      test-session
      my-worktree
      """
    And session "my-worktree" contains following windows:
      | window_name       |
      | WorktreeInherited |
