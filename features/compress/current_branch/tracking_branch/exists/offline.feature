Feature: compress the commits in offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And offline mode is enabled
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
    And an uncommitted file
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git add -A                  |
      |         | git stash -m "Git Town WIP" |
      |         | git reset --soft main       |
      |         | git commit -m "commit 1"    |
      |         | git stash pop               |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE  |
      | feature | local    | commit 1 |
      |         | origin   | commit 1 |
      |         |          | commit 2 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                               |
      | feature | git add -A                            |
      |         | git stash -m "Git Town WIP"           |
      |         | git reset --hard {{ sha 'commit 2' }} |
      |         | git stash pop                         |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now
    And the uncommitted file still exists
