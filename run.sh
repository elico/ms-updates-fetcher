#!/usr/bin/env bash

lockfile -r 0 /tmp/store-fetcher.lock || exit 1

#/usr/bin/store-fetcher -dir=/var/storedata
/usr/bin/store-fetcher -dir=/home/eliezer/Scripts/requeststore/storedata/

echo ""
echo ""
echo "removing lock file: /tmp/store-fetcher.lock"
rm -f /tmp/store-fetcher.lock

