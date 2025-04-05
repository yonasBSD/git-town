Feature: prototype the current perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot prototype perennial branches
      """
    And the perennial branches are still "perennial"
    And there are still no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the perennial branches are still "perennial"
    And there are still no prototype branches
