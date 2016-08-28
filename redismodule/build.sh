#!/bin/bash

[ -f redismodule.h ] || wget https://raw.githubusercontent.com/antirez/redis/unstable/src/redismodule.h
gcc -fPIC -std=gnu99 -c -o reporter.o reporter.c
gcc -fPIC -std=gnu99 -c -o cmd_router.o cmd_router.c
ld -o cmd_router.so cmd_router.o reporter.o -shared -Bsymbolic -lc -lffi -lnanomsg
