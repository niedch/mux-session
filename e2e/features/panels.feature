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
      window_name = "HorizontalSplit"
      
      [[project.window.panel_config]]
      panel_direction = "h"
      
      [[project.window.panel_config]]
      panel_direction = "h"

      [[project.window]]
      window_name = "VerticalSplit"
      
      [[project.window.panel_config]]
      panel_direction = "v"
      
      [[project.window.panel_config]]
      panel_direction = "v"
      """
    Then I expect following sessions:
      """
      test-session
      my-project
      """
    And session "my-project" contains following windows:
      | window_name     |
      | HorizontalSplit |
    And window "HorizontalSplit" in session "my-project" is split vertically
    And window "VerticalSplit" in session "my-project" is split horizontally
