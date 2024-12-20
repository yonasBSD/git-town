Feature: display the parent of a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | perennial | perennial |        | local, origin |
    And the current branch is "perennial"
    When I run "git-town config get-parent"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints no output
