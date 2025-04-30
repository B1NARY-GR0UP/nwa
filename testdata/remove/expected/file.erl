-module(hello).
-export([hello_world/0]).

hello_world() ->
    io:format("Hello, World!~n").

% A simple function with pattern matching
greet(Name) ->
    io:format("Hello, ~s!~n", [Name]).