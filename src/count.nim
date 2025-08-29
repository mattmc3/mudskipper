import posix

# Exported function: counts positional arguments plus number of stdin lines if stdin not a TTY.
proc count*(args: seq[string]): int =
  var n = args.len
  # posix.isatty returns a cint; compare to 0 instead of using 'not'.
  if isatty(0) == 0:
    while true:
      try:
        discard stdin.readLine()
        inc n
      except EOFError:
        break
  return n

when isMainModule:
  import os
  echo count(commandLineParams())
