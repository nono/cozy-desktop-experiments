#!/usr/bin/env bash
set -e

rm -rf tmp* workspace/
mkdir workspace
cd desktop
(
	sleep 1
	cd ../workspace
	mkdir foo
	echo bar > bar
	sleep 20
	mv bar foo/
) &
node index.js
