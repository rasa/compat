#!/usr/bin/env sh

# to run script locally
: "${GITHUB_WORKSPACE:=${PWD}}"
: "${GOOS:=$(uname | tr '[:upper:]' '[:lower:]')}" || true
: "${GOARCH:=$(uname -p)}" || true
case "${GOARCH}" in
  x86_64)
    GOARCH=amd64
    ;;
  *) ;;
esac

if ! command -v gtar >/dev/null 2>/dev/null; then
  gtar() { tar "$@"; }
fi

if ! command -v sha256sum >/dev/null 2>/dev/null; then
  sha256sum() { gsha256sum "$@"; }
fi

printf 'GOOS:   %s\n' "${GOOS}"
printf 'GOARCH: %s\n' "${GOARCH}"

tmp1=$(mktemp)
curl -L -s -o "${tmp1}" 'https://go.dev/dl/?mode=json'
jqcmd="[ .[] | select(.stable == true) ][0] | .files[] | select(.os == \"${GOOS}\" and .arch == \"${GOARCH}\")"

unset GOOS GOARCH

name=$(jq -r "${jqcmd} | .filename" "${tmp1}")
printf 'name:   %s\n' "${name}"
hash=$(jq -r "${jqcmd} | .sha256" "${tmp1}")
printf 'hash:   %s\n' "${hash}"
size=$(jq -r "${jqcmd} | .size" "${tmp1}")
printf 'size:   %s\n' "${size}"
base=$(basename "${name}" .tar.gz)
printf 'base:   %s\n' "${base}"

mkdir -p "../${base}"
cd "../${base}" || exit

printf 'Downloading %s...\n' "https://go.dev/dl/${name}"
curl -L -s -o "${name}" "https://go.dev/dl/${name}"

printf '%s %s\n' "${hash}" "${name}" | sha256sum -c

printf 'Untarring %s to %s...\n' "${name}" "${PWD}"
gtar xzf "${name}"

rm -f "${name}" "${tmp1}"

export PATH="${PWD}/go/bin:${PATH}"

cd "${GITHUB_WORKSPACE}" || exit

GOVERSION=$(go version || true)
printf 'gover:  %s\n' "${GOVERSION}"

rv=0

# NOTE: dragonflybsd requires -buildvcs=false
if go build -buildvcs=false -trimpath ./...; then
  printf '::notice ::build succeeded: %s\n' "${GOVERSION}"
else
  rv=$?
  printf '::error ::build failed: %s (error %s)\n' "${GOVERSION}" "${rv}"
fi

if go test .; then
  printf '::notice ::tests succeeded: %s\n' "${GOVERSION}"
else
  rv=$?
  printf '::error ::tests failed: %s (error %s)\n' "${GOVERSION}" "${rv}"
fi

exit "${rv}"
