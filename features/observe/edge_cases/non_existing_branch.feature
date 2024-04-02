Feature: cannot observe non-existing branches

  Background:
    Given the current branch is a feature branch "feature"
    And an uncommitted file
    When I run "git-town observe feature non-existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is still "feature"
    And the uncommitted file still exists
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are still no observed branches
    And the current branch is still "feature"