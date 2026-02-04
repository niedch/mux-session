Feature: Config validation
  As a user
  I want to ensure that my config is correctly configured
  So that mux-session can operate correctly

  Scenario: Setting Primary Window multiple times result in error
    Given a new tmux server
    And I have the following directories:
      | name       |
      | my-project |
    When I run mux-session config-validate with config:
      """
      search_paths = ["<search_path>"]
      
      [[project]]
      name = "my-project"
      
      [[project.window]]
      window_name = "Main"
      primary = true
      
      [[project.window]]
      window_name = "Sub"
      primary = true
      """
    Then I should see the following items in output:
      | lines                                                             |
      | only one window can be marked as primary in project configuration |
