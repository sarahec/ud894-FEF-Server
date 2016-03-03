
#!/bin/bash

PLATFORMS="linux_amd64 darwin_amd64 windows_amd64"
rm -rf binaries &>/dev/null; mkdir binaries

find . -maxdepth 1 -mindepth 1 -type d \
	| grep -v '^\./\.' \
	| while read dir; do
	echo $dir
	for platform in $PLATFORMS; do
		export GOARCH=${platform##*_}
		export GOOS=${platform%%_*}
		go build -o binaries/${dir}_${platform} ${dir}
	done
done
