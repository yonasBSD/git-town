# Commands

Run `git town` for an overview of all Git Town commands and
`git town help <command>` for help with individual commands. Each Git Town
command can be called like `git town <command>`. This user manual displays the
commands in the shorter form available after running
[git town alias](commands/alias.md).

### Typical development commands

- [git hack](commands/hack.md) - create a new feature branch
- [git sync](commands/sync.md) - update the current branch with all ongoing
  changes
- [git new-pull-request](commands/new-pull-request.md) - create a new pull
  request
- [git ship](commands/ship.md) - deliver a completed feature branch

### Advanced development commands

- [git kill](commands/kill.md) - delete a feature branch
- [git prune-branches](commands/prune-branches.md) - remove all merged branches
- [git rename-branch](commands/rename-branch.md) - rename a branch
- [git repo](commands/repo.md) - view the Git repository in the browser

### Nested feature branches

- [git append](commands/append.md) - create a new feature branch as a child of
  the current branch
- [git prepend](commands/prepend.md) - create a new feature branch between the
  current branch and its parent
- [git town set-parent-branch](commands/set-parent-branch.md) - change the
  parent of a feature branch

### Git Town configuration

- [git town config](commands/config.md) - display or update your Git Town
  configuration
- [git town new-branch-push-flag](commands/new-branch-push-flag.md) - configure
  whether new empty branches are pushed to origin
- [git town main-branch](commands/main-branch.md) - display/set the main
  development branch for the current repo
- [git town offline](commands/offline.md) - enable/disable offline mode
- [git town perennial-branches](commands/perennial-branches.md) - display or
  update the perennial branches for the current repo
- [git town pull-branch-strategy](commands/pull-branch-strategy.md) - display or
  set the strategy to update perennial branches

### Git Town setup

- [git town alias](commands/alias.md) - add or remove shorter aliases for Git
  Town commands
- [git town completion](commands/completion.md) - generate completion scripts
  for Bash, zsh, fish & PowerShell.
- [git town version](commands/version.md) - display the installed version of Git
  Town