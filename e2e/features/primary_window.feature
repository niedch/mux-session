Feature: Primary Window Selection
  As a user
  I want to define a primary window
  So that it is focused when I open the session

  Scenario: Default behavior focuses the last created window
    Given a new tmux server
    And I have the following directories:
      | name       |
      | my-project |
    When I run mux-session switch "my-project" with config:
      """
      search_paths = ["<search_path>"]
      
      [[project]]
      name = "my-project"
      
      [[project.window]]
      window_name = "First"

      [[project.window]]
      window_name = "Second"
      """
    Then session "my-project" contains following windows:
      | window_name |
      | First       |
      | Second      |
    And window "Second" in session "my-project" is active

  Scenario: Primary flag focuses the specified window
    Given a new tmux server
    And I have the following directories:
      | name       |
      | my-project |
    When I run mux-session switch "my-project" with config:
      """
      search_paths = ["<search_path>"]
      
      [[project]]
      name = "my-project"
      
      [[project.window]]
      window_name = "First"
      primary = true

      [[project.window]]
      window_name = "Second"
      """
    Then session "my-project" contains following windows:
      | window_name |
      | First       |
      | Second      |
    And window "First" in session "my-project" is active
