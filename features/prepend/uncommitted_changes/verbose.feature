Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And an uncommitted file

  Scenario: result
    When I run "git-town prepend parent --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                       |
      |        | backend  | git version                                   |
      |        | backend  | git rev-parse --show-toplevel                 |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git status --long --ignore-submodules         |
      |        | backend  | git stash list                                |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}     |
      |        | backend  | git remote                                    |
      |        | backend  | git remote get-url origin                     |
      |        | backend  | git log main..old --format=%s --reverse       |
      | old    | frontend | git add -A                                    |
      |        | frontend | git stash -m "Git Town WIP"                   |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      | old    | frontend | git checkout -b parent main                   |
      |        | backend  | git show-ref --verify --quiet refs/heads/main |
      |        | backend  | git config git-town-branch.parent.parent main |
      |        | backend  | git show-ref --verify --quiet refs/heads/old  |
      |        | backend  | git config git-town-branch.old.parent parent  |
      |        | backend  | git show-ref --verify --quiet refs/heads/old  |
      |        | backend  | git stash list                                |
      | parent | frontend | git stash pop                                 |
      |        | backend  | git branch -vva --sort=refname                |
      |        | backend  | git config -lz --includes --global            |
      |        | backend  | git config -lz --includes --local             |
      |        | backend  | git stash list                                |
    And Git Town prints:
      """
      Ran 27 shell commands.
      """
    And the current branch is now "parent"

  Scenario: undo
    Given I ran "git-town prepend parent"
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                          |
      |        | backend  | git version                                      |
      |        | backend  | git rev-parse --show-toplevel                    |
      |        | backend  | git config -lz --includes --global               |
      |        | backend  | git config -lz --includes --local                |
      |        | backend  | git status --long --ignore-submodules            |
      |        | backend  | git stash list                                   |
      |        | backend  | git branch -vva --sort=refname                   |
      |        | backend  | git remote get-url origin                        |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}        |
      |        | backend  | git remote get-url origin                        |
      | parent | frontend | git add -A                                       |
      |        | frontend | git stash -m "Git Town WIP"                      |
      |        | frontend | git checkout old                                 |
      | old    | frontend | git branch -D parent                             |
      |        | backend  | git config --unset git-town-branch.parent.parent |
      |        | backend  | git config git-town-branch.old.parent main       |
      |        | backend  | git stash list                                   |
      | old    | frontend | git stash pop                                    |
    And Git Town prints:
      """
      Ran 18 shell commands.
      """
    And the current branch is now "old"
