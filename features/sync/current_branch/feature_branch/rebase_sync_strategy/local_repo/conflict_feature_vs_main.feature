@skipWindows
Feature: handle conflicts between the current feature branch and the main branch (in a local repo)

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                          |
      | feature | git rebase main --no-update-refs |
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
    And the current branch is still "feature"
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND            |
      | feature | git rebase --abort |
    And the current branch is still "feature"
    And no rebase is now in progress
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "feature"
    And a rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and enter "resolved commit" for the commit message
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                   |
      | feature | git -c core.editor=true rebase --continue |
    And all branches are now synchronized
    And the current branch is still "feature"
    And no rebase is now in progress
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and enter "resolved commit" for the commit message
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND |
