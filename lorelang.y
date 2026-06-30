%{
#include <stdio.h>
#include <stdlib.h>

extern int yylex();
extern int yylineno;
void yyerror(const char *s);
%}

/* Declaración de Tokens */
%token TOKEN_PERSONAJE TOKEN_ATRIBUTOS TOKEN_RESTRICCIONES TOKEN_ALRECIBIR
%token TOKEN_INYECTAR TOKEN_SI TOKEN_SINO TOKEN_RESPONDER TOKEN_GENERAR TOKEN_FUNCION
%token ID STRING

%%

/* Reglas Gramaticales de LoreLang */

origen : personaje_bloque { printf("[Sintactico] Estructura completa de LoreLang valida.\n"); }
       ;

personaje_bloque : TOKEN_PERSONAJE STRING '{' secciones '}' 
                 { printf("[Oracion Valida] Declaracion de Personaje exitosa.\n"); }
                 ;

secciones : bloque_atributos_declaracion bloque_restricciones_declaracion bloques_evento funcion_bloque
          ;

/* Atributos */
bloque_atributos_declaracion : TOKEN_ATRIBUTOS '{' '}' lista_asignaciones
                             | TOKEN_ATRIBUTOS '{' lista_asignaciones '}'
                             | lista_asignaciones
                             ;

lista_asignaciones : lista_asignaciones asignacion
                   | /* vacio */
                   ;

asignacion : ID ':' STRING ';' { printf("[Oracion Valida] Atributo/Propiedad declarada.\n"); }
           ;

/* Restricciones */
bloque_restricciones_declaracion : TOKEN_RESTRICCIONES '{' '}' lista_restricciones
                                 | TOKEN_RESTRICCIONES '{' lista_restricciones '}'
                                 | lista_restricciones
                                 ;

lista_restricciones : lista_restricciones restriccion
                    | /* vacio */
                    ;

restriccion : ID ':' '[' lista_strings ']' ';' { printf("[Oracion Valida] Restriccion de lista indexada.\n"); }
            | ID ':' STRING ';'                { printf("[Oracion Valida] Restriccion escalar declarada.\n"); }
            ;

lista_strings : STRING
              | lista_strings ',' STRING
              ;

/* Eventos */
bloques_evento : bloques_evento evento
               | /* vacio */
               ;

evento : TOKEN_ALRECIBIR '(' STRING ')' '{' TOKEN_INYECTAR ':' STRING ';' '}' 
       { printf("[Oracion Valida] Disparador de Evento 'AlRecibir' configurado correctamente.\n"); }
       ;

/* Funciones */
funcion_bloque : TOKEN_FUNCION ID '(' ID ')' '{' cuerpo_funcion '}'
               { printf("[Oracion Valida] Bloque de Funcion de procesamiento reconocido.\n"); }
               ;

cuerpo_funcion : condicional_si
               ;

condicional_si : TOKEN_SI '(' ID ID STRING ')' '{' instruccion_interna '}' bloque_sino
               { printf("[Oracion Valida] Estructura de control condicional 'Si' procesada.\n"); }
               ;

bloque_sino : TOKEN_SINO '{' instruccion_interna '}'
            | /* vacio */
            ;

instruccion_interna : TOKEN_RESPONDER '(' STRING ')' ';'
                    | TOKEN_GENERAR '(' ID ')' ';'
                    ;

%%

void yyerror(const char *s) {
    fprintf(stderr, "Error sintactico en la linea %d: %s\n", yylineno, s);
}

int main(int argc, char **argv) {
    printf("Iniciando Analisis Sintactico de LoreLang...\n");
    return yyparse();
}