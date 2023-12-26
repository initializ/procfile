#!/usr/bin/env bash

set -euo pipefail

GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/initializ/procfile/v5/cmd/main
GOOS="windows" GOARCH="amd64" go build -ldflags='-s -w' -o bin/main.exe github.com/initializ/procfile/v5/cmd/main

if [ "${STRIP:-false}" != "false" ]; then
  strip bin/main bin/main.exe
fi

if [ "${COMPRESS:-none}" != "none" ]; then
  $COMPRESS bin/main bin/main.exe
fi

ln -fs main bin/build
ln -fs main bin/detect
ln -fs main.exe bin/build.exe
ln -fs main.exe bin/detect.exe
