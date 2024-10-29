Feature: does not compress an active observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE    | FILE NAME  | FILE CONTENT |
      | observed | local, origin | observed 1 | observed_1 | observed 1   |
      |          |               | observed 2 | observed_2 | observed 2   |
    And the branches
      | NAME  | TYPE    | PARENT   | LOCATIONS     |
      | child | feature | observed | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | child  | local, origin | child 1 | child_1   | child 1      |
      |        |               | child 2 | child_2   | child 2      |
    And the current branch is "observed"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                         |
      | observed | git fetch --prune --tags                        |
      |          | git add -A                                      |
      |          | git stash                                       |
      |          | git checkout child                              |
      | child    | git reset --soft observed                       |
      |          | git commit -m "child 1"                         |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout observed                           |
      | observed | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "observed"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE    |
      | child    | local, origin | child 1    |
      | observed | local, origin | observed 1 |
      |          |               | observed 2 |
    And file "observed_1" still has content "observed 1"
    And file "observed_2" still has content "observed 2"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                         |
      | observed | git add -A                                      |
      |          | git stash                                       |
      |          | git checkout child                              |
      | child    | git reset --hard {{ sha 'child 2' }}            |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout observed                           |
      | observed | git stash pop                                   |
    And the current branch is still "observed"
    And the initial commits exist now
    And the initial branches and lineage exist now
    And the uncommitted file still exists
