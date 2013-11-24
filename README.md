nogo
====

nogo is a command line utility that will create and parse a directory of
markdown files for simplify organized note taking.

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
