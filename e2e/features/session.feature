Feature: Mux Session functionality
  As a user of mux-session
  I want to manage my terminal sessions
  So that I can work efficiently

  Background:
    Given I build the mux-session

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
