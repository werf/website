#!/bin/bash

INPUT_THROTTLING="${INPUT_THROTTLING:-11}"

while read p; do
	if [ "${INSTANT_TYPE}x" != "x" ] ; then
		delay=${TYPE_DELAY:-0.1}
		echo -n '$ ' ; echo "${p}" ; sleep $delay ; eval ${p} ; sleep $delay
	else
		echo -n '$ ' ; echo "${p}" | pv -qL ${INPUT_THROTTLING} ; eval ${p}
	fi
done < ${1}
