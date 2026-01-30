Feature: Mux Session functionality
  As a user of mux-session
  I want to manage my terminal sessions
  So that I can work efficiently

  Scenario: Basic Help view
    Given I build the mux-session
    When I run mux-session with help flag
    Then I should see help output
  
  Scenario: Basic session creation
    Given a new tmux server
    Then I expect following Sessions:
    """
    doc string test
    """
    
