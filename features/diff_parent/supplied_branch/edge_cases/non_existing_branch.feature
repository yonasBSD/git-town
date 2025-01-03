Feature: does not diff non-existing branch

  Scenario: non-existing branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town diff-parent non-existing"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      there is no branch "non-existing"
      """
