Feature: rename a branch to itself

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"

  Scenario: without force
    When I run "git-town rename old"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot rename branch to current name
      """

  Scenario: with force
    When I run "git-town rename --force old"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot rename branch to current name
      """
