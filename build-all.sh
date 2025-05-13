#!/usr/bin/env bash

set +e

mapfile -t targets < <(go tool dist list | grep -E -v '(android|ios)' || true)

declare -A seen
rv=0
for target in "${targets[@]}"; do
  export GOOS="${target%%/*}"
  if [[ -v seen[${GOOS}] ]]; then
    continue
  fi
  seen[${GOOS}]=1
  export GOARCH="${target#*/}"
  echo "Building for ${GOOS}/${GOARCH}"
  go build -v .
  ((rv |= $?))
done
exit "${rv}"
