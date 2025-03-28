@skipWindows
Feature: display all executed Git commands

  Scenario:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --verbose"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                   |
      |        | backend  | git version                               |
      |        | backend  | git rev-parse --show-toplevel             |
      |        | backend  | git config -lz --includes --global        |
      |        | backend  | git config -lz --includes --local         |
      |        | backend  | which wsl-open                            |
      |        | backend  | which garcon-url-handler                  |
      |        | backend  | which xdg-open                            |
      |        | backend  | which open                                |
      |        | backend  | git branch --show-current                 |
      | (none) | frontend | open https://github.com/git-town/git-town |
    And Git Town prints:
      """
      Ran 10 shell commands.
      """
