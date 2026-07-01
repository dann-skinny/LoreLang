lex lorelang.lex
yacc -d lorelang.y
gcc y.tab.c lex.yy.c -lfl -o lorelang
./lorelang < tom_nook.lore
