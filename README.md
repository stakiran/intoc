# intoc

TOC generator for Markdown.

<!-- toc -->
- [intoc](#intoc)
  - [Feature](#feature)
  - [How to use](#how-to-use)
    - [Usage](#usage)
  - [Samples](#samples)
    - [Basic](#basic)
    - [Depth control](#depth-control)
    - [Use not hyphen but asterisk](#use-not-hyphen-but-asterisk)
    - [Direct update](#direct-update)
    - [Plain enumeration](#plain-enumeration)
    - [No link format but keep list grammer](#no-link-format-but-keep-list-grammer)
  - [How to develop](#how-to-develop)
    - [Requirement](#requirement)
    - [How to run](#how-to-run)
    - [How to test](#how-to-test)
    - [How to build](#how-to-build)
  - [License](#license)
  - [Author](#author)

:warning: This intoc is new version based on golang. The old python based version is [here](py).

## Feature

- Golang based.
- No WebAPI use.
- Multiple output ways.
  - Output to the stdout
  - Direct Update(Inserting TOC to the next of `<!-- toc -->` line directly).
- Support sections written in Japanese.

## How to use

```
$ intoc -input (Target-Markdown-File)
```

Create an alias if needed.

### Usage

```
$ intoc -h
Usage of intoc:
  -debug-print-all
        [DEBUG] print all options with name and value.
  -edit
        If given then insert TOC to the file from '-input'.
  -edit-target string
        A insertion destination label when '-edit' given. NOT CASE-SENSITIVE. (default "<!-- TOC")
  -indent-depth int
        The number of spaces per a nest in TOC. (default 2)
  -input string
        A input filename.
  -no-linkformat
        Not use '[text](#anochor)', but use 'text'.
  -parse-depth int
        The depth of the TOC list nesting. If minus then no limit depth. (default -1)
  -use-asterisk
        Use an asterisk '*' as a list grammer.
  -use-plain-enum
        Not use Markdown grammer, but use simple plain section name listing.
  -version
        Show this intoc version.
```


## Samples

### Basic

`-input` is required option.

```
$ intoc -input readme.md
- [intoc](#intoc)
  - [Feature](#feature)
  - [How to use](#how-to-use)
    - [Usage](#usage)
  - [Samples](#samples)
    - [Basic](#basic)
    - [Depth control](#depth-control)
    - [Use not hyphen but aasterisk](#use-not-hyphen-but-aasterisk)
    - [Direct update](#direct-update)
    - [Plain enumeration](#plain-enumeration)
    - [No link format but keep list grammer](#no-link-format-but-keep-list-grammer)
  - [How to develop](#how-to-develop)
    - [Requirement](#requirement)
    - [How to run](#how-to-run)
    - [How to test](#how-to-test)
    - [How to build](#how-to-build)
  - [License](#license)
  - [Author](#author)

```

### Depth control

`-indent-depth` and `-parse-depth`.

```
$ intoc -input readme.md -indent-depth 4 -parse-depth 2
- [intoc](#intoc)
    - [Feature](#feature)
    - [How to use](#how-to-use)
    - [Samples](#samples)
    - [How to develop](#how-to-develop)
    - [License](#license)
    - [Author](#author)
```

### Use not hyphen but asterisk

`-use-asterisk`.

```
$ intoc -input readme.md -use-asterisk
* [intoc](#intoc)
  * [Feature](#feature)
  * [How to use](#how-to-use)
    * [Usage](#usage)
  * [Samples](#samples)
    * [Basic](#basic)
    * [Depth control](#depth-control)
    * [Use not hyphen but aasterisk](#use-not-hyphen-but-aasterisk)
    * [Direct update](#direct-update)
    * [Plain enumeration](#plain-enumeration)
    * [No link format but keep list grammer](#no-link-format-but-keep-list-grammer)

  * [How to develop](#how-to-develop)
    * [Requirement](#requirement)
    * [How to run](#how-to-run)
    * [How to test](#how-to-test)
    * [How to build](#how-to-build)
  * [License](#license)
  * [Author](#author)
```

### Direct update

Write `<!-- toc -->` to your input file and use `-edit`.

```
$ intoc -input readme.md -edit

$ type readme.md
# intoc

TOC generator for Markdown.

<!-- toc -->
- [intoc](#intoc)
  - [Feature](#feature)
  - [How to use](#how-to-use)
    - [Usage](#usage)
  - [Samples](#samples)
    - [Basic](#basic)
    - [Depth control](#depth-control)
    - [Use not hyphen but asterisk](#use-not-hyphen-but-asterisk)
    - [Direct update](#direct-update)
    - [Plain enumeration](#plain-enumeration)
    - [No link format but keep list grammer](#no-link-format-but-keep-list-grammer)
  - [How to develop](#how-to-develop)
    - [Requirement](#requirement)
    - [How to run](#how-to-run)
    - [How to test](#how-to-test)
    - [How to build](#how-to-build)
  - [License](#license)
  - [Author](#author)

## Feature
...
```

If you want to change edit-string from default `<!-- toc -->` to an another, use `-edit-target` option.

### Plain enumeration

`-use-plain-enum`.

```
$ intoc -input readme.md -use-plain-enum
intoc
Feature
How to use
Usage
Samples
Basic
Depth control
Use not hyphen but asterisk
Direct update
Plain enumeration
No link format but keep list grammer
How to develop
Requirement
How to run
How to test
How to build
License
Author
```

### No link format but keep list grammer

`-no-linkformat`.

```
$ intoc -input readme.md -no-linkformat
- intoc
  - Feature
  - How to use
    - Usage
  - Samples
    - Basic
    - Depth control
    - Use not hyphen but asterisk
    - Direct update
    - Plain enumeration
    - No link format but keep list grammer
  - How to develop
    - Requirement
    - How to run
    - How to test
    - How to build
  - License
  - Author
```

## How to develop

### Requirement

- Golang
  - Developed and Tested on `go version go1.10.3 windows/386` and `go version go1.10.3 windows/amd64`
- github.com/stretchr/testify/assert for unittest

### How to run

See [run.bat](run.bat), but executing `go run` simply.

### How to test

See [test.bat](test.bat), but executing `go test` simply.

### How to build

See [build.bat](build.bat), but executing `go build` simply.

## License

[MIT License](LICENSE)

## Author

[stakiran](https://github.com/stakiran)
