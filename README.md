nogo
====

nogo is a command line utility that will create and parse a directory of
markdown files for simplified organized note taking. nogo will maintain a
simple directory structure and will launch your editor of choice to edit your
notes.

## Examples

```
~ nogo

nogo - the notes helper

actions:

  nogo help

  nogo new
  nogo new [topic]
  nogo new [topic] [event]

  nogo ls
  nogo ls [topic]

  nogo edit
  nogo edit [topic]
  nogo edit [topic] [note name substring]

~ nogo ls

looks like there aren't any topcis to list!

~ nogo new
Enter the notes topic: this
Enter the event name: that
~ nogo ls

all topics:
   this

~ nogo ls this

notes in this:
   that (2013-11-23)

~
```

### Completion and substrings

Most `nogo` commands don't actually require you to write out the whole topic
or event name. For example, if you have a topic "meetings" and an event
"meeting with my manager about cats", you can type `nogo edit mee cats` and
`nogo` will automatically find the note you're looking for.

## Configuration

### Editor

By default, this is set to `vim`. If you don't like `vim`, you can set the
`$EDITOR` environment variable to your test based editor of choice. You likely
already have this set (to `vim` or `emacs` for example) if you use the command
line to do things that require you to edit files.

### Notes path

By default, this is set to be `~/notes`. Simply set the `NOGODIR` environment
variable if you'd like this to be different.

```
export NOGODIR="/Users/marpaia/Desktop"
```

## Requirements

The only requirement is that you have Go installed. If you don't have Go
installed and, for some reason, would like to keep it that way, let me know
and I'll upload a nogo binary somewhere for you to download.

You can install go from http://golang.org/doc/install or via a package manager.

## Installation

Use `go get` to install nogo:

```
go get github.com/marpaia/nogo
```

## External dependencies

This project has no external dependencies other than the Go standard library.

## Hacking

If you'd like to customize nogo, do so as you would any other Go project. Use
`go get` as described above to download and install the software, edit the code
in whatever way you see fit, and execute `go install github.com/marpaia/nogo`.

## Contributing

Please contribute and help improve this project!

- Fork the repo
- Improve the code
- Submit a pull request

## Areas for improvement

Check out the issues labeled "enhancement" for easy contribution/improvement
ideas. This might be especially interesting to you if you'd like to cut your
teeth contributing to open source software as this is a pretty simple piece of
code that can be improved in a lot of ways.

https://github.com/marpaia/nogo/issues?labels=enhancement
