#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

if [[ "$arg_site_lang" == "en" ]]; then
  language="en_US"
elif [[ "$arg_site_lang" == "ru" ]]; then
  language="ru_RU"
fi

for file in `find ./ -type f -name "*.html"`
do
   wc -l $file;
   stat -c %s $file;
done
