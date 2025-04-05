Feature: create a branch in the presence of a symbolic reference

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And I ran "git symbolic-ref refs/remotes/origin/master refs/remotes/origin/main"
    And the current branch is "main"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git fetch --prune --tags                |
      |        | git rebase origin/main --no-update-refs |
      |        | git checkout -b new                     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | new    | git checkout main                           |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch -D new                           |
    And the initial commits exist now
    And no lineage exists now
