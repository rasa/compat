#!/usr/bin/env bash
# SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT

set +e

if test -n "${BUILD_ALL:-}"; then
  ignore='zzz'
else
  ignore='x(386|arm)$'
fi
mapfile -t targets < <(go tool dist list | grep -E -v '(android|ios)/' | grep -E -v "${ignore}" || true)

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
  go test -c "$@" ./...
  ((rv |= $?))
  if ((rv > 0)); then
    exit "${rv}"
  fi
done
exit "${rv}"
