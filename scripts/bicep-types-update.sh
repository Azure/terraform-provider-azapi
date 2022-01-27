#!/bin/bash

MYDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOTDIR="$MYDIR/.."

main() {
    bicep_dir=$1
    cp -r "$bicep_dir/generated" "$ROOTDIR/internal/azure/generated"
    cd $ROOTDIR/internal/azure/generated
    find . -name "*.md" -type f -delete
    find . -name "*.out" -type f -delete
}

main "$@"