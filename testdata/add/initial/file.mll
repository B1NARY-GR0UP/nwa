{
open Parser
}

rule token = parse
  | [' ' '\t' '\n']    { token lexbuf }
  | "hello"            { HELLO }
  | "world"            { WORLD }
  | eof                { EOF }