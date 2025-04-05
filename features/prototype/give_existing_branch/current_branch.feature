Feature: prototype the current feature branch

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town prototype feature"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "feature" is now a prototype branch
      """
    And the prototype branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And there are now no prototype branches
    And the initial branches and lineage exist now
