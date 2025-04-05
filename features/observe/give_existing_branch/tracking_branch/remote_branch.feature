Feature: make another remote feature branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE    | PARENT | LOCATIONS |
      | remote-feature | feature | main   | origin    |
    And I run "git fetch"
    When I run "git-town observe remote-feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      |        | git checkout remote-feature |
    And Git Town prints:
      """
      branch "remote-feature" is now an observed branch
      """
    And branch "remote-feature" now has type "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH         | COMMAND                      |
      | remote-feature | git checkout main            |
      | main           | git branch -D remote-feature |
    And there are now no observed branches
