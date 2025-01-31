Feature: displaying the branches in the middle of an ongoing sync merge conflict

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    And I ran "git-town sync" and ignore the error
    When I run "git-town branch"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   feature
      """
