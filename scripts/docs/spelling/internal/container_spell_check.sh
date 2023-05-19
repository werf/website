#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"
str=$'\n'

if [[ "$arg_site_lang" == "en" ]]; then
  language="en_US"
elif [[ "$arg_site_lang" == "ru" ]]; then
  language="ru_RU,en_US"
fi

html2text 404.html | hunspell -d $language -l

#for file in `find ./ -type f -name "*.html"`
#do
#  echo "$str"
#  echo "Checking $file..."
#  cat $file
#  html2text $file | hunspell -d $language -l
#done
