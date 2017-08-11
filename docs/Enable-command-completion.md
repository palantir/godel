gödel supports command completion on Bash and Zsh using the functionality provided by the
[`pkg/palantir/cli` package](https://github.com/palantir/pkg/tree/develop/cli). Because gödel is almost always invoked
using the wrapper script `./godelw`, it is recommended that completion be set up on `godelw`.

# Bash

1. Run `./godelw _completion --prog=godelw --alias ./godelw --bash` in a project that has `godel`.
2. This command will output text to the terminal:
```
_completion_######() {
...
}
compdef _completion_###### "godelw"
# Run 'eval "$(/<path...>/.godel/dists/godel-<version>/bin/darwin-amd64/godel --wrapper /<path...>/godelw _completion --prog=godelw --alias ./godelw --bash)"' to apply
```
3. Copy all of the text except for the comment (starting at `_completion_######() {` and ending at `compdef _completion_###### "godelw"`).
4. Add the command to the shell's startup environment (`~/.bash_profile` or equivalent):
```
_completion_######() {
...
}
compdef _completion_###### "godelw"
```

# Zsh

1. Run `./godelw _completion --prog=godelw --alias ./godelw --zsh` in a project that has the `godel` version that should
   be used for completion.
2. This command will output text to the terminal:
```
_completion_######() {
...
}
compdef _completion_###### "godelw"
# Run 'eval "$(/<path...>/.godel/dists/godel-<version>/bin/darwin-amd64/godel --wrapper /<path...>/godelw _completion --prog=godelw --alias ./godelw --zsh)"' to apply
```
3. Copy all of the text except for the comment (starting at `_completion_######() {` and ending at `compdef _completion_###### "godelw"`).
4. Add the command to the shell's startup environment (`~/.zshrc` or equivalent):
```
_completion_######() {
...
}
compdef _completion_###### "godelw"
```

# Update Completion
The steps above should continue to work as long as the content of the completion script itself (or the way it is
invoked) does not change. Such a change should not be common, and if it does occur it will be noted in the release
notes. If the completion script does change, it can be updated by repeating the steps above using the new version of
`godelw`.

# Explanation
The `--prog` flag determines the command for which completion should be performed. Because completion should be
performed on `./godelw`, the value of `--prog` is set to `godelw` (Bash and Zsh both handle command completion on
relative path invocation automatically). The `--alias` flag specifies the program that should be invoked to generate the
completion. Because it is expected that completion will be performed in a directory that contains `./godelw`, the value
is specified as `./godelw`. This ensures that command completion is provided by the proper version of gödel that is
being used (which may be important in environments in which multiple versions of gödel exist).

This setup assumes that `godelw` will always be invoked using `./godelw`. If this is not the case and completion is
desired in those scenarios as well, the value provided go `--alias` can be changed to be the absolute path to a `godel`
executable, which will make it such that the specified executable always generates the completions for the command. This
has the down-side of tying all completions of `./godelw` to a specific version, but does allow for completion on
invocations from any location.
