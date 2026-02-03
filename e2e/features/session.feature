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
