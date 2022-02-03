@skipWindows
Feature: ask for missing configuration

  Scenario:
    Given I haven't configured Git Town yet
    When I run "git-town kill" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
    And my repo now has no perennial branches
    And it prints the error:
      """
      you can only kill feature branches
      """
    And Git Town still has no branch hierarchy information
