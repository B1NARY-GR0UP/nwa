file { '/tmp/hello.txt':
  ensure  => file,
  content => "Hello, World!\n",
  mode    => '0644',
}

# Define a class
class hello {
  notify { 'hello_world':
    message => 'Hello, World from Puppet!',
  }
}

# Include the class
include hello