#!/bin/bash

INPUT_THROTTLING="${INPUT_THROTTLING:-13}"

i=0
while read p; do
	if [[ "${INSTANT_TYPE}x" != "x" ]] || [ "${i}" -eq "0" ] ; then
		delay=${TYPE_DELAY:-0.1}
		echo -n '$ ' ; echo "${p}" ; sleep $delay ; eval ${p}
	else
		echo -n '$ ' ; echo "${p}" | pv -qL ${INPUT_THROTTLING} ; eval ${p}
	fi

	((i++))
done < ${1}
