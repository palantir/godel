Name
====
gödel's name is an homage to [Kurt Gödel](https://en.wikipedia.org/wiki/Kurt_G%C3%B6del). As with many Go tools, it is a
play on words that involves "go". It is also a play on [Gödel's incompleteness theorems](https://en.wikipedia.org/wiki/G%C3%B6del's_incompleteness_theorems),
with the idea being that the standard Go tools provide a consistent system for Go projects, but there are truths that
cannot be proven to be true using just the standard Go tooling itself. gödel acts as a tool outside of this system.

Usage of name
-------------
The name of this project is gödel. The 'g' is lowercase, and the second letter is the precomposed representation of an
'o' with a diaeresis (`ö`, `U+00F6`). This representation is used in the documentation and in the source code.

However, for any file that can be written to disk, "godel" is used as the name instead. This is done to preserve maximal
compatibility with all terminals and for filesystems that do not fully support Unicode names. The GitHub project name is
also "godel" because GitHub only supports the character class `[A-Za-z0-9_.-]` for project names.

Sidebar: HFS+ and decomposed unicode
------------------------------------
At the beginning of the project, an attempt was made to use "gödel" as the canonical name for everything (including
files and settings on-disk). Unfortunately, this was complicated by the fact that HFS+ (which is the default file system
used by most MacOS systems) normalizes all Unicode names using Normalization Form D (NFD).

Quick background: Unicode has a notion of [precomposed or decomposed characters](https://en.wikipedia.org/wiki/Precomposed_character).
Unicode supports representing the `ö` character in two different ways: a precomposed form (`ö`, `U+00F6`) and a
decomposed form ('o' + '¨', `U+006F` + `U+0308`). Unicode does define a notion of
[equivalence](https://en.wikipedia.org/wiki/Unicode_equivalence#Combining_and_precomposed_characters), and Unicode
equivalence recognizes these two forms as equivalent. However, by default, most calls compare strings strictly as byte
sequences, in which case `ö` and `o¨` are distinct.

Most file systems write file names exactly as they are provided. However, this is not the case for HFS+: HFS+ explicitly
normalizes all file names using [Normalization Form D](http://unicode.org/reports/tr15/#Norm_Forms), which performs
"canonical decomposition" on all inputs.

Thus, when a file named `ö` (`U+00F6`) is written to disk, HFS+ translates it into `o¨` (`U+006F` + `U+0308`). The
system-level calls handles translation, so a request to open `ö` is translated into a request for `o¨`. However, from a
data perspective, this means that reads and writes can be asymmetric: requests for files with precomposed Unicode
characters will return names that contain decomposed characters.

This poses problems when trying to be interoperable with other systems that do not have this restriction. For example,
in most file systems it is perfectly legal to have one file named `ö` and another named `o¨`. However, this is simply
not possible to represent in HFS+. This is similar to the case insensitivity restriction -- by default, the HFS+ file
system is case-insensitive, meaning that it is not possible to have files in a directory that differ only based on the
case of the characters in the names (however, it is possible to format an HFS+ volume to be case-sensitive, which is
recommended for most developers). These are known quirks of HFS+ that cause numerous headaches (some of which are
enumerated in an oft-cited [rant by Linus Torvalds](https://plus.google.com/+JunioCHamano/posts/1Bpaj3e3Rru)). Luckily,
it appears that this issue has been fixed correctly in the [Apple File System (APFS)](https://en.wikipedia.org/wiki/Apple_File_System).

Initially, an attempt was made to normalize the names in gödel code to deal with this. In Go, this can be done by
importing the `golang.org/x/text/unicode/norm` package and using `norm.NFC.String` to convert Unicode strings into
NFC-normalized format. This would ensure that any instances of `go¨del` would be converted to `gödel`.

In Bash scripts, the following function was used to return `go¨del` on Darwin systems and `gödel` on all other systems:

```
normalize() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # HFS file systems use deconstructed UTF-8
        echo $(iconv -t UTF8-MAC <<<$1)
    else
        echo $1
    fi
}
```

Although these work-arounds functioned correctly, in writing them it was clear that this approach would not be
sustainable. It isn't realistic to require all gödel users to know about Unicode normalization and deal with it in their
own tooling, and there seemed to be glitches in every new external system that had to interact with these characters
(for example, URL-encoding of the characters was also handled differently by different systems when uploading artifacts
that contained the name).
