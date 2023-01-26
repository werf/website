#!/bin/bash

echo -n '$ ' ; echo "$@" | pv -qL 10 ; $@
