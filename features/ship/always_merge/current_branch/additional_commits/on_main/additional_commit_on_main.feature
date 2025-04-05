Feature: can ship not-up-to-date feature branches using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | main    | local, origin | main commit    |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                             |
      | feature | git fetch --prune --tags            |
      |         | git checkout main                   |
      | main    | git merge --no-ff --edit -- feature |
      |         | git push                            |
      |         | git push origin :feature            |
      |         | git branch -D feature               |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | main commit            |
      |        |               | feature commit         |
      |        |               | Merge branch 'feature' |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And the currently checked out commit is "feature commit"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | main commit            |
      |        |               | feature commit         |
      |        |               | Merge branch 'feature' |
    And the initial branches and lineage exist now
