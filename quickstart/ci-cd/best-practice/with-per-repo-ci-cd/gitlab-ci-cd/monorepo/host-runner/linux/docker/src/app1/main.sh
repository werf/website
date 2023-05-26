#!/bin/sh

while true; do
    printf "HTTP/1.1 200 OK\n\nHello world.\n" | ncat -lp 80
done
