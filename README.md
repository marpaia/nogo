nogo
====

nogo is a command line utility that will create and parse a directory of
markdown files for simplified organized note taking. nogo will maintain a
simple directory structure and will launch your editor of choice to edit your
notes.

## Configuration

### Editor

Set the `$EDITOR` environment variable to your test based editor of choice. You
likely already have this set (to `vim` for example) if you use the command line
to do things that require you to edit files.

### Notes path

This is hardcoded at `~/notes` in the `nogo.go` file. Simply change the
`notesSubDir` variables from `notes` to whatever you want if you'd like your
notes to be saved somewhere else. The final path is just a concatenation of
`"$HOME/{notesSubDir}"`.

## Requirements

The only requirement is that you have Go installed. If you don't have Go
installed and, for some reason, would like to keep it that way, let me know
and I'll upload a nogo binary somewhere for you to download.

## Installation

Use `go get` to install nogo:
```
go get github.com/marpaia/nogo
```

## External dependencies

This project has no external dependencies other than the Go standard library.

## Examples

```
~ nogo

nogo - the notes helper

actions:
  nogo new
  nogo new [topic]
  nogo new [topic] [event]
  nogo ls
  nogo ls [topic]
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
