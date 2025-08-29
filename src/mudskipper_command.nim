# Base command abstractions for Mudskipper

# Swappable writers export
type Writer* = proc (s: string)

# Base command providing writer hooks
type MudskipperCommand* = ref object of RootObj
  outWriter*: Writer
  errWriter*: Writer

# Abstract-like usage method: must be overridden by subclasses
method usage*(self: MudskipperCommand) {.base.} =
  raise newException(CatchableError, "usage() not implemented")

proc output*(self: MudskipperCommand, s: string) =
  if self.outWriter == nil: stdout.write(s) else: self.outWriter(s)

proc outerr*(self: MudskipperCommand, s: string) =
  if self.errWriter == nil: stderr.write(s) else: self.errWriter(s)
