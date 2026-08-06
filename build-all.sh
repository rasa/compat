#!/usr/bin/env bash
# SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT

set +e

mapfile -t targets < <(go tool dist list | grep -E -v '(android|ios)/' || true)

declare -A seen
rv=0
for target in "${targets[@]}"; do
  export GOOS="${target%%/*}"
  if [[ -v seen[${GOOS}] ]]; then
    continue
  fi
  test -n "${BUILD_ALL:-}" || seen[${GOOS}]=1
  export GOARCH="${target#*/}"
  echo "*** Building for ${GOOS}/${GOARCH}: build args: $*"
  go build -v "$@" .
  ((rv |= $?))
  if ((rv>0)); then
    exit "${rv}"
  fi
done
exit "${rv}"
