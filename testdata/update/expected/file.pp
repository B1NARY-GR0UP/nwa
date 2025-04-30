# Copyright 2025 BINARY Members
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
