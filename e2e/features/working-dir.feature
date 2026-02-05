Feature: Mux Session Creating a new Session

  Scenario: Creating Session with working dir of the choosen item
    Given a new tmux server
    And I have the following directories:
      | name          |
      | project-one   |
      | project-two   |
      | project-three |
      | my-project    |
    When I run mux-session switch "my-project" with config:
      """
      search_paths = ["<search_path>"]
      
      [[default.window]]
      window_name = "Shell"
      cmd = "pwd"
      """
    Then I expect following sessions:
      """
      test-session
      my-project
      """
    When I execute following Command in Session "my-project" on Window "Shell":
      """
      pwd
      """
    Then I should see the following lines in output:
      | lines      |
      | my-project |
