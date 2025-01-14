# git town continue

```command-summary
git town continue [-v | --verbose]
```

When a Git Town command encounters a problem that it cannot resolve, for example
a merge conflict, it stops to give the user an opportunity to resolve the issue.
Once you have resolved the issue, run the _continue_ command to tell Git Town to
continue executing the failed command. Git Town will retry the failed operation
and execute all remaining operations of the original command.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
