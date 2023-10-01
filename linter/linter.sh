#!/bin/sh

ROOT=$(dirname "$(echo "$0" | grep -E "^/" -q && echo "$0" || echo "$PWD/${0#./}")")

info() {
  echo "\033[0;32m[INFO][linter] $*\033[0m"
}
err() {
  echo "\033[0;31m[ERR][linter] $*\033[0m"
}

MSG="Linter will be launched"
echo "$ARGS" | grep '\--new' -q && MSG="${MSG} only for new files" || MSG="${MSG} for all files"

info "$MSG"

which "golangci-lint" 1>/dev/null
if [ $? -eq 0 ]; then
  info "Linter will start with the installed application"

  # shellcheck disable=SC2068
  golangci-lint run --timeout 3m \
    --config "${ROOT}/golangci.yml" \
    --color always \
    --out-format colored-line-number \
    --fast $@

  exit $?
fi

which "docker" 1>/dev/null
if [ $? -eq 0 ]; then
  info "Linter will run in a docker container"

  # shellcheck disable=SC2068
  docker run --rm \
    -w /app \
    -v "${ROOT}/../:/app:ro" \
    golangci/golangci-lint:v1.54.2 \
    golangci-lint run -v \
    --new \
    --config /app/linter/golangci.yml \
    --timeout 3m \
    --color always \
    --out-format colored-line-number $@

  exit $?
fi

err "To run the linter, you need to install docker or golangci-lint"
exit 1