#!/usr/bin/env ruby
# Hello World in Ruby

puts "Hello, World!"

# Define a class
class Greeter
  def initialize(name)
    @name = name
  end

  def greet
    puts "Hello, #{@name}!"
  end
end

# Create an instance and call method
greeter = Greeter.new("Ruby")
greeter.greet