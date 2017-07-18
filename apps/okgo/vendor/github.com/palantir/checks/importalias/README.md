importalias
===========
`importalias` is a check that verifies that import aliases in a project are used consistently. It verifies that, if a
given package is imported using an alias, then all aliases for that import are consistent across the project.

Usage
-----
`importalias` uses its current working directory as the project root. If no arguments are provided, it is invoked on all
of the go packages it can find in the current working directory and its subdirectories. If arguments are provided, they
are interpreted as packages relative to the working directory, and only the specified packages will be checked.

By default, the output of the check is standard Go check output format. The program operates as follows:

* Finds all imports and all aliases that are used for imports.
* If a package is imported using multiple different aliases, the alias that is most commonly used to import the package
  is considered the "correct" import.
  * If there is a tie for the most commonly used alias, it is assumed that there is no consensus for the alias.
* Any line that imports a package using an alias that is not the most common one (or an alias for which there is no
  consensus) is treated as an error. The file and line number is printed, along with a suggestion for how the alias
  should be renamed.

The `-v` or `--verbose` flag can be used to print an overview of all of the imports in the project that are imported
using multiple aliases. The output is organized by import and lists all of the aliases used for the import (in order of
most commonly used) and the files and locations in the files in which the imports occur.
