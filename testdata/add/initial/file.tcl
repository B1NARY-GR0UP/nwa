#!/usr/bin/env tclsh
# Hello World in Tcl

puts "Hello, World!"

# Define a procedure
proc greet {name} {
    puts "Hello, $name!"
    return $name
}

# Call the procedure
greet "Tcl"