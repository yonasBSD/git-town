Feature: handle conflicts between the current observed branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the current branch is "observed"
    And the commits
      | BRANCH   | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | observed | local    | conflicting local commit  | conflicting_file | local content  |
      |          | origin   | conflicting origin commit | conflicting_file | origin content |
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                     |
      | observed | git fetch --prune --tags                    |
      |          | git add -A                                  |
      |          | git stash -m "Git Town WIP"                 |
      |          | git rebase origin/observed --no-update-refs |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """
    And a rebase is now in progress
    And the uncommitted file is stashed

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND            |
      | observed | git rebase --abort |
      |          | git stash pop      |
    And the current branch is still "observed"
    And the uncommitted file still exists
    And no rebase is now in progress
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the uncommitted file is stashed
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                   |
      | observed | git -c core.editor=true rebase --continue |
      |          | git stash pop                             |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                   |
      | observed | local, origin | conflicting origin commit |
      |          | local         | conflicting local commit  |
    And the current branch is still "observed"
    And no rebase is now in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH   | NAME             | CONTENT          |
      | observed | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND       |
      | observed | git stash pop |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                   |
      | observed | local, origin | conflicting origin commit |
      |          | local         | conflicting local commit  |
    And the current branch is still "observed"
    And no rebase is now in progress
    And the uncommitted file still exists
    And these committed files exist now
      | BRANCH   | NAME             | CONTENT          |
      | observed | conflicting_file | resolved content |
