Feature: Environment Variable Management
  As a user
  I want to configure environment variables for my projects
  So that my sessions have the correct environment settings

  Scenario: Project config sets environment variables
    Given a new tmux server
    And I have the following directories:
      | name       |
      | my-project |
    When I run mux-session switch "my-project" with config:
      """
      search_paths = ["<search_path>"]
      
      [[project]]
      name = "my-project"
      
      [project.env]
      MY_PROJECT_VAR = "production_value"
        
      [[project.window]]
      window_name = "Main"
      
      [[project.window]]
      window_name = "Sub"
      """
    Then I expect following sessions:
      """
      test-session
      my-project
      """
    And session "my-project" contains following windows:
      | window_name |
      | Main        |
      | Sub         |
    When I execute following Command in Session "my-project" on Window "Main":
      """
      export | grep MY_PROJECT_VAR
      """
    Then I should see the following items in output:
      | item                                  |
      | MY_PROJECT_VAR=("?production_value"?) |
