Feature: auto-creating a prototype branch when hacking

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And Git Town setting "create-prototype-branches" is "true"
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                 |
      | existing | git fetch --prune --tags                |
      |          | git checkout main                       |
      | main     | git rebase origin/main --no-update-refs |
      |          | git checkout -b new                     |
    And the current branch is now "new"
    And branch "new" is now prototype
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | local, origin | main commit     |
      | existing | local         | existing commit |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                     |
      | new      | git checkout main                           |
      | main     | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout existing                       |
      | existing | git branch -D new                           |
    And the current branch is now "existing"
    And the initial commits exist now
    And the initial branches and lineage exist now
