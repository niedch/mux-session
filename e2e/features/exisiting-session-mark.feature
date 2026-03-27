Feature: Mux Session mark existing Sessions
  As a user of mux-session
  I want to see which sessions are already their
  So that I can work efficiently

  Scenario: When i have a session already running i want this selection to be prefix with "[x]"
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
      | item               |
      |  test-session     |
      | 󰄱 .*/project-one   |
      | 󰄱 .*/project-three |
      | 󰄱 .*/project-two   |
    When I run mux-session switch "project-one" with config:
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
      | item               |
      |  test-session     |
      | [] .*/project-one |
      | 󰄱 .*/project-three |
      | 󰄱 .*/project-two   |

  Scenario: Tmux internal Sessions should not be marked with "[ ]"
    Given a new tmux server
    And I have the following directories:
      | name        |
      | project-one |
    When I run mux-session list-sessions with config:
      """
      search_paths = ["<search_path>"]
      """
    Then I should see the following items in output:
      | item           |
      |  test-session |
