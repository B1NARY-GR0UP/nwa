#!/usr/bin/env perl
# Hello World in Perl

use strict;
use warnings;

print "Hello, World!\n";

# Define a subroutine
sub greet {
    my $name = shift;
    print "Hello, $name!\n";
}

# Call the subroutine
greet("Perl");