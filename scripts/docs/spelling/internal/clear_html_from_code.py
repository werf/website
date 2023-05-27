#!/usr/bin/python3
# -*- coding: utf-8 -*-

from bs4 import BeautifulSoup
import sys


if len (sys.argv) > 1:
    try:
        html = open(sys.argv[1]).read()
        root = BeautifulSoup(html, 'html.parser')
        try:
            for code in root.select('code'):
                code.decompose()
            print(root)
        except:
            print("Code tag not found.")
    except:
        print("Filename error.")
