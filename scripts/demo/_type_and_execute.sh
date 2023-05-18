#!/bin/bash

INPUT_THROTTLING="${INPUT_THROTTLING:-11}"

while read p; do
	if [ "${INSTANT_TYPE}x" == "1x" ] ; then
		echo -n '$ ' ; echo "${p}" ; sleep 0.8 ; eval ${p}
	else
		echo -n '$ ' ; echo "${p}" | pv -qL ${INPUT_THROTTLING} ; eval ${p}
	fi
done < ${1}
