#!/usr/bin/env sh
# ~/download-build-test.sh
# Cownload go, build, and test code
# Called by the github actions test-*bsd.yml

set -vx

# to run script locally
: "${GITHUB_REPOSITORY:=rasa/$(basename "${PWD}")}"
: "${GITHUB_WORKSPACE:=${PWD}}"
: "${GOOS:=$(uname | tr '[:upper:]' '[:lower:]')}" || true
: "${GOARCH:=$(uname -p)}" || true
: "${GOOPTS:=}"
: "${GOVERSION:=1.25.0}"
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

if [ $# -gt 0 ]; then
  GOOPTS="$@"
fi

env | sort

printf 'CODECOV_SLUG:      %s\n' "${CODECOV_SLUG:-}"
# shellcheck disable=SC2154 # (warning): CODECOV_TOKEN is referenced but not assigned.
printf 'CODECOV_TOKEN:     %d chars long\n' "${#CODECOV_TOKEN}"
printf 'GITHUB_REPOSITORY: %s\n' "${GITHUB_REPOSITORY}"
printf 'GITHUB_WORKSPACE:  %s\n' "${GITHUB_WORKSPACE}"
printf 'GOARCH:            %s\n' "${GOARCH}"
printf 'GOOS:              %s\n' "${GOOS}"
printf '$GOOPTS:           %s\n' "${GOOPTS}"

tmp1=$(mktemp)
curl -L -s -o "${tmp1}" 'https://go.dev/dl/?mode=json'
jqcmd="[ .[] | select(.stable == true) ][0] | .files[] | select(.os == \"${GOOS}\" and .arch == \"${GOARCH}\")"

# gover=$(jq -r " ${jqcmd} | .[0].version | ltrimstr("go")")
# printf 'gover:  %s\n' "${gover}"

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
printf 'GOVERSION:         %s\n' "${GOVERSION}"

# NOTE: dragonflybsd requires -buildvcs=false
printf "Running: "
echo go build -buildvcs=false -trimpath $GOOPTS ./...
if ! go build -buildvcs=false -trimpath $GOOPTS ./...; then
  rv=$?
  printf '::error ::build failed: %s (error %s)\n' "${GOVERSION}" "${rv}"
  exit "${rv}"
fi

printf '::notice ::build succeeded: %s\n' "${GOVERSION}"

printf "Running: "
echo go test -covermode=atomic -coverprofile=coverage.out -coverpkg=. $GOOPTS -v .
if ! go test -covermode=atomic -coverprofile=coverage.out -coverpkg=. $GOOPTS -v .; then
  rv=$?
  printf '::error ::tests failed: %s (error %s)\n' "${GOVERSION}" "${rv}"
  exit "${rv}"
fi

printf '::notice ::tests succeeded: %s\n' "${GOVERSION}"

sed -i.bak "/compat\/cmd\//d; /compat\/golang\//d;" coverage.out
rm -f coverage.out.bak

# ls -l

exit 0
