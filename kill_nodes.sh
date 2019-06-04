#!/usr/bin/env bash
lsof -ti:10001 | xargs kill -9
lsof -ti:10002 | xargs kill -9
lsof -ti:10003 | xargs kill -9
lsof -ti:10004 | xargs kill -9
rm *.csv