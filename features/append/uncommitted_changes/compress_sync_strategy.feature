Feature: append a new feature branch in a dirty workspace using the "compress" sync strategy

# TODO: move to a subfolder of "clean_workspace"
# and remove the "clean_workspace" folder from the folder hierarchy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |
      | existing | local, origin | existing commit 2 |
    And the current branch is "existing"
    And an uncommitted file
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout -b new         |
      | new      | git stash pop               |
      |          | git restore --staged .      |
    And the current branch is now "new"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | new      | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout existing       |
      | existing | git branch -D new           |
      |          | git stash pop               |
      |          | git restore --staged .      |
    And the current branch is now "existing"
    And the initial commits exist now
    And the initial lineage exists now