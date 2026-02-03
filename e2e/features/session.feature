Feature: Mux Session functionality
  As a user of mux-session
  I want to manage my terminal sessions
  So that I can work efficiently

  Scenario: Basic Help view
    When I run mux-session with help flag
    Then I should see help output

  Scenario: Basic session creation
    Given a new tmux server
    Then I expect following sessions:
    """
    test-session
    """

  Scenario: List sessions shows directories from search path
    Given a new tmux server
    And I have the following directories:
      | name          |
      | project-one   |
      | project-two   |
      | project-three |
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    Then I should see the following items in output:
      | item                |
      | project-one         |
      | project-two         |
      | project-three       |
      | [TMUX] test-session |

  Scenario: Switch command creates new session from directory
    Given a new tmux server
    And I have the following directories:
      | name       |
      | my-project |
    When I run mux-session switch "my-project" with config:
      """
      search_paths = ["<search_path>"]

      [default]
      [[default.window]]
      window_name = "Shell"
      """
    Then I expect following sessions:
      """
      test-session
      my-project
      """
    And session "my-project" contains following windows:
      | window_name |
      | Shell       |
