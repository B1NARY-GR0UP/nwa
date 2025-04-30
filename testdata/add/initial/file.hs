module Main where

main :: IO ()
main = putStrLn "Hello, World!"

-- A simple function
greet :: String -> String
greet name = "Hello, " ++ name ++ "!"