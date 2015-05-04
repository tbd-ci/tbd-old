#!/usr/bin/env bash
set -eo pipefail
IFS=$'\n\t'

CI_DIR=${CI_DIR:-ci}

set -u

WORKTREE=$( tbd-save_tree )
NOW=$( date +%Y%m%d%H%M%S )
HOST=$( hostname )

# -tree defaults to current if not specified
BUILDS=$( tbd-lookup-builds \
  -tree $( $WORKTREE ) \
  -target spec:covered \
  -exit-code 0 \
)

if [ -z "$BUILDS" ] ; then
  >&2 echo "No builds found"
  exit 1
fi

for TASK in $( ls -1 "$CI_DIR" )
do
  if [ ! -e "$TASK/dependencies" ] ; then
    echo $TASK
  else
    local UNMET_DEPS=""

    for DEPENDENCY in $( cat "$TASK/dependencies" )
    do
      local BUILDS
      # Find successful builds of this target/tree combination
      BUILDS=$( tbd-lookup-builds \
        -tree "$WORKTREE" \
        -target "$DEPENDENCY" \
        -exit-code 0 \
      )
      if [ -z "$BUILDS" ] ; then
        UNMET_DEPS="$UNMET_DEPS\n$DEPENDENCY"
      fi
    done
    if [ -z "UNMET_DEPS" ] ; then
      # All dependencies are met
      echo $TASK
    fi
  fi
done
