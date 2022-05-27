#!/bin/bash

MYDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOTDIR="$MYDIR/.."

main() {
    bicep_dir=$1
    echo "$bicep_dir/generated"
    echo "$ROOTDIR/internal/azure/generated"
    echo "removing all exist type files..."
    rm -r "$ROOTDIR/internal/azure/generated"
    echo "done"
    echo "copying new type files"
    cp -r "$bicep_dir/generated" "$ROOTDIR/internal/azure/generated"
    echo "done"
    cd $ROOTDIR/internal/azure/generated
    find . -name "*.md" -type f -delete
    find . -name "*.out" -type f -delete
}

main "$@"