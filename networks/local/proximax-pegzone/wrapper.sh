#!/usr/bin/env sh

BINARY=/pxbd/${BINARY:-pxbd}
ID=${ID:-0}
LOG=${LOG:-pxbd.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'pxbd' E.g.: -e BINARY=pxbd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export PXBDHOME="/pxbd/node${ID}/pxbd"

if [ -d "$(dirname "${PXBDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${PXBDHOME}" "$@" | tee "${PXBDHOME}/${LOG}"
else
  "${BINARY}" --home "${PXBDHOME}" "$@"
fi

