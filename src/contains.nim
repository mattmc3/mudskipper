import os

proc printHelp() =
  echo """
Usage: contains [-i|--index|-0|--index0] NEEDLE HAYSTACK...
  -i, --index      Print the 1-based index of NEEDLE in HAYSTACK
  -0, --index0     Print the 0-based index of NEEDLE in HAYSTACK
  -h, --help       Show this help message
Returns 0 if NEEDLE is found, 1 otherwise.
"""

proc contains*(args: seq[string]): int =
  if args.len == 0 or args[0] == "-h" or args[0] == "--help":
    printHelp()
    return 0

  var indexFlag = false
  var indexStart = 1
  var remainingArgs = args

  # Shift off flag if present
  if remainingArgs.len > 0 and (remainingArgs[0] == "-i" or remainingArgs[0] == "--index"):
    indexFlag = true
    indexStart = 1
    remainingArgs = remainingArgs[1..^1]
  elif remainingArgs.len > 0 and (remainingArgs[0] == "-0" or remainingArgs[0] == "--index0"):
    indexFlag = true
    indexStart = 0
    remainingArgs = remainingArgs[1..^1]

  if remainingArgs.len < 1:
    echo "contains: Key not specified"
    return 2

  let needle = remainingArgs[0]
  let haystack = remainingArgs[1..^1]

  for i, item in haystack:
    if item == needle:
      if indexFlag:
        echo $(i + indexStart)
      return 0

  return 1

when isMainModule:
  import os
  quit(contains(commandLineParams()))
