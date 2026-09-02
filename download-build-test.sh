#!/usr/bin/env sh
# SPDX-FileCopyrightText: Copyright (c) 2026 Ross Smith II <ross@smithii.com>
# SPDX-License-Identifier: MIT
# ~/download-build-test.sh
# Download go, build, and test code
# Called by the github actions test-{*bsd|illumos|solaris}.yml

# to run script locally
goos=$(uname | tr '[:upper:]' '[:lower:]' || true)
goarch=$(uname -m || true)
goversion=$(grep '^go [1-9]\.' go.mod | cut -d' ' -f 2 || true)

: "${GITHUB_REPOSITORY:=rasa/$(basename "${PWD}")}"
: "${GITHUB_WORKSPACE:=${PWD}}"
: "${GOOS:=${goos}}"
: "${GOARCH:=${goarch}}"
: "${GOOPTS:=}"
: "${GOVERSION:=${goversion}}"

case "${GOARCH}" in
  386|amd64|arm|arm64|loong64|mips|mips64|mips64le|mipsle|ppc64|ppc64le|riscv64|s390x)
    ;;
  x86_64)
    GOARCH=amd64
    ;;
  aarch64)
    GOARCH=arm64
    ;;
  arm*)
    GOARCH=arm
    ;;
  *86)
    GOARCH=386
    ;;
  *)
    printf 'Unsupported CPU architecture: %s\n' "${GOARCH}" >&2
    exit 1
    ;;
esac

if ! command -v gtar >/dev/null 2>/dev/null; then
  gtar() { tar "$@"; }
fi

if ! command -v sha256sum >/dev/null 2>/dev/null; then
  sha256sum() { gsha256sum "$@"; }
fi

if [ $# -gt 0 ]; then
  GOOPTS="$*"
fi

printf 'CODECOV_SLUG:      %s\n' "${CODECOV_SLUG:-}"
# shellcheck disable=SC2154 # (warning): CODECOV_TOKEN is referenced but not assigned.
printf 'CODECOV_TOKEN:     %d chars long\n' "${#CODECOV_TOKEN}"
printf 'GITHUB_REPOSITORY: %s\n' "${GITHUB_REPOSITORY}"
printf 'GITHUB_WORKSPACE:  %s\n' "${GITHUB_WORKSPACE}"
printf 'GOARCH:            %s\n' "${GOARCH}"
printf 'GOOS:              %s\n' "${GOOS}"
printf 'GOOPTS:            %s\n' "${GOOPTS}"
printf 'GOVERSION:         %s\n' "${GOVERSION}"

tmp1=$(mktemp)
url='https://go.dev/dl/?mode=json&include=all'
printf 'Downloading %s...\n' "${url}"
curl -L -s -o "${tmp1}" "${url}"

if [ -n "${GOVERSION}" ]; then
  # shellcheck disable=SC2016
  jqcmd='.[] | select(.version == $version) | .files[] | select(.os == $os and .arch == $arch  and .kind == "archive")'
else
  # shellcheck disable=SC2016
  jqcmd='[ .[] | select(.stable == true) ][0] | .files[] | select(.os == $os and .arch == $arch)'
fi

name=$(jq --arg os "${GOOS}" --arg arch "${GOARCH}" --arg version "go${GOVERSION}" -r "${jqcmd} | .filename" "${tmp1}")
printf 'name:   %s\n' "${name}"

hash=$(jq --arg os "${GOOS}" --arg arch "${GOARCH}" --arg version "go${GOVERSION}" -r "${jqcmd} | .sha256" "${tmp1}")
printf 'hash:   %s\n' "${hash}"

size=$(jq --arg os "${GOOS}" --arg arch "${GOARCH}" --arg version "go${GOVERSION}" -r "${jqcmd} | .size" "${tmp1}")
printf 'size:   %s\n' "${size}"

base=$(basename "${name}" .tar.gz)
printf 'base:   %s\n' "${base}"

cd "${GITHUB_WORKSPACE}/.." || exit

godir="${PWD}/${base}"
if [ ! -d "${godir}" ]; then
  printf 'Creating directory %s\n' "${godir}"
  mkdir -p "${godir}"
  cd "${godir}" || exit

  url="https://go.dev/dl/${name}"
  printf 'Downloading %s to %s...\n' "${url}" "${name}"
  curl -L -s -o "${name}" "${url}"

  printf '%s %s\n' "${hash}" "${name}" | sha256sum -c

  printf 'Untarring %s to %s...\n' "${name}" "${PWD}"
  gtar xzf "${name}"

  rm -f "${name}" "${tmp1}"
fi

cd "${GITHUB_WORKSPACE}" || exit

export PATH="${godir}/go/bin:${PATH}"

GOVER=$(go version || true)
printf 'Go ver: %s\n' "${GOVER}"

# NOTE: dragonflybsd requires -buildvcs=false
printf "Running: "
# shellcheck disable=SC2086 # (info): Double quote to prevent globbing and word splitting.
echo go build -buildvcs=false -trimpath ${GOOPTS} ./...
rv=0
# shellcheck disable=SC2086 # (info): Double quote to prevent globbing and word splitting.
go build -buildvcs=false -trimpath ${GOOPTS} ./... || rv=$?
if [ "${rv}" != "0" ]; then
  printf '::error ::build failed: %s (error %s)\n' "${GOVER}" "${rv}"
  exit "${rv}"
fi

printf '::notice ::build succeeded: %s\n' "${GOVER}"

printf "Running: "
# shellcheck disable=SC2086 # (info): Double quote to prevent globbing and word splitting.
echo go test -covermode=atomic -coverprofile=coverage.out -coverpkg=. ${GOOPTS} -v ./...
rv=0
# shellcheck disable=SC2086 # (info): Double quote to prevent globbing and word splitting.
go test -covermode=atomic -coverprofile=coverage.out -coverpkg=. ${GOOPTS} -v ./... || rv=$?
if [ "${rv}" != "0" ]; then
  printf '::error ::tests failed: %s (error %s)\n' "${GOVER}" "${rv}"
  exit "${rv}"
fi

printf '::notice ::tests succeeded: %s\n' "${GOVER}"

sed -i.bak "/compat\/cmd\//d; /compat\/golang\//d;" coverage.out
rm -f coverage.out.bak

exit 0
