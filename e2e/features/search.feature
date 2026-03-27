Feature: Search functionality
  As a user of mux-session
  I want to filter items with a search query
  So that I can find specific sessions quickly

  Scenario: List sessions with search filter
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
    And I search for "project-two"
    Then I should see the following items in output:
      | item        |
      | project-two |
    And I should not see "project-one" in output
    And I should not see "project-three" in output

  Scenario: List sessions with partial search
    Given a new tmux server
    And I have the following directories:
      | name          |
      | project-alpha |
      | project-beta  |
      | random-dir    |
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    And I search for "project"
    Then I should see the following items in output:
      | item          |
      | project-alpha |
      | project-beta  |
    And I should not see "random-dir" in output

  Scenario: List sessions with multiple search results
    Given a new tmux server
    And I have the following directories:
      | name          |
      | project       |
      | project-alpha |
      | project-beta  |
      | random-dir    |
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    And I search for "project"
    Then I should see the following items in output:
      | item          |
      | project       |
      | project-alpha |
      | project-beta  |
    And I should not see "random-dir" in output
