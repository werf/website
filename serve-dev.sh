#!/bin/bash

source $(trdl use werf 1.2 ea)

werf compose up --follow --dev --docker-compose-command-options="-d" jekyll_base
