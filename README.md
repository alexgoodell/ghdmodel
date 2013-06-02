# Global Health Decisions - Data Models

The data models underlying GHD and a simple web server to make them available over the network.

## Install

```bash
cd $GOPATH
mkdir -p github.com/alexgoodell
git clone git@github.com:alexgoodell/ghd-modeling.git github.com/alexgoodell/ghimodel
```

## Use

```bash
go install
ghdmodel
```

## Debug

Use gdb as debugger.

OSX users need to codesign their gdb, see here: http://sourceware.org/gdb/wiki/BuildingOnDarwin

General information on debugging go programs: http://golang.org/doc/gdb

```bash
go install -gcflags "-N -l"
gdb ghdmodel
```


Copyright © 2013 Alex Goodell