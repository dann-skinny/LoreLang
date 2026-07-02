#! /usr/bin/env bash

lex lorelang.lex

yacc -d lorelang.y

cc y.tab.c lex.yy.c -lfl -o lorelang

./lorelang < tom_nook.lore
