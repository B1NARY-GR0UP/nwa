%{
open Ast
%}

%token HELLO WORLD EOF
%start <unit> prog

%%

prog:
  | HELLO WORLD EOF { print_endline "Hello, World!" }
;