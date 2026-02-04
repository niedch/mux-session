Feature: Panel Configuration
  As a user
  I want to configure panels in my windows
  So that I can have multiple panes with different commands

  Scenario: Switch command creates new session with multiple panels and commands
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
      window_name = "MultiPanel"

      [[project.window.panel_config]]
      panel_direction = "h"
      cmd = "cat"

      [[project.window.panel_config]]
      panel_direction = "h"
      cmd = "tail -f /dev/null"

      [[project.window.panel_config]]
      panel_direction = "v"
      cmd = "sleep 5"
      """
    Then I expect following sessions:
      """
      test-session
      my-project
      """
    And session "my-project" contains following windows:
      | window_name |
      | MultiPanel  |
    And window "MultiPanel" in session "my-project" contains following panels:
      | command |
      | cat     |
      | tail    |
      | sleep   |
