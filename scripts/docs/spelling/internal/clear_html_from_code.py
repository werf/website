#!/usr/bin/python3
# -*- coding: utf-8 -*-

from bs4 import BeautifulSoup
import sys


if len (sys.argv) > 1:
    html = open(sys.argv[1]).read()

    root = BeautifulSoup(html, 'html.parser')
    body = root.select_one('body')
    for code in body.select('code'):
        code.decompose()

    print(root)