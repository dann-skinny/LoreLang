%{
#include <stdio.h>
#include <stdlib.h>
/* Importamos la definición de tokens generada por Yacc automáticamente */
#include "y.tab.h"
%}

%option noyywrap
%option yylineno

%%

[ \t\n\r]+              { /* Ignorar espacios, tabuladores y saltos de línea */ }
"//".* { /* Ignorar comentarios de una línea */ }

"Personaje"             { return TOKEN_PERSONAJE; }
"Atributos"             { return TOKEN_ATRIBUTOS; }
"Restricciones"         { return TOKEN_RESTRICCIONES; }
"AlRecibir"             { return TOKEN_ALRECIBIR; }
"InyectarContexto"      { return TOKEN_INYECTAR; }
"Si"                    { return TOKEN_SI; }
"Sino"                  { return TOKEN_SINO; }
"contiene"              { return TOKEN_CONTIENE; }
"Responder"             { return TOKEN_RESPONDER; }
"GenerarRespuesta"      { return TOKEN_GENERAR; }
"Funcion"               { return TOKEN_FUNCION; }

[a-zA-Z_][a-zA-Z0-9_]* { return ID; }
\"([^\\\"]|\\.)*\"      { return STRING; }

"{"                     { return '{'; }
"}"                     { return '}'; }
"("                     { return '('; }
")"                     { return ')'; }
"["                     { return '['; }
"]"                     { return ']'; }
":"                     { return ':'; }
";"                     { return ';'; }
","                     { return ','; }

.                       { printf("Error lexico en la linea %d: Caracter no reconocido '%s'\n", yylineno, yytext); }

%%