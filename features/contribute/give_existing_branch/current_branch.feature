Feature: make the current contribution branch a contribution branch

  Background:
    Given a local Git repo
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | contribution | contribution | main   | local     |
    When I run "git-town contribute contribution"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "contribution" is already a contribution branch
      """
    And the contribution branches are still "contribution"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the contribution branches are still "contribution"
