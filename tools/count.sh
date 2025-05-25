#!/usr/bin/env dash
#/bin/sh
# Optionally, #!/usr/bin/env dash

##? count - counts args

(( REPLY=$# ))

# count piped/redirected input
if ! [ -t 0 ]; then
  while IFS= read -r DATA || [ -n "$DATA" ]; do
    (( REPLY+=1 ))
  done
fi
printf '%s\n' $REPLY
