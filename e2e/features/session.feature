Feature: Mux Session functionality
  As a user of mux-session
  I want to manage my terminal sessions
  So that I can work efficiently

  Scenario: Basic session creation
    Given I have mux-session installed
    When I run mux-session with help flag
    Then I should see help output