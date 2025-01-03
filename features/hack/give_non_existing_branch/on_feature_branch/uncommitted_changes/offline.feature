Feature: offline mode

  Background:
    Given a Git repo with origin
    And offline mode is enabled
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | main   | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout -b new         |
      | new    | git stash pop               |
    And the current branch is now "new"
    And the uncommitted file still exists
    And the initial commits exist now
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | new    | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout main           |
      | main   | git branch -D new           |
      |        | git stash pop               |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the initial commits exist now
    And no lineage exists now
