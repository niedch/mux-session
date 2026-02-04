Feature: Mux Session help
  As a user of mux-session
  I want to see the Help menu of mux-session

  Scenario: Basic Help view
    When I run mux-session with help flag
    Then I should see help output
