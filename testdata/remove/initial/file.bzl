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

def _hello_impl(ctx):
    output = ctx.outputs.out
    ctx.actions.write(
        output = output,
        content = "Hello, World!\n",
    )

# Define a rule
hello = rule(
    implementation = _hello_impl,
    attrs = {
        "message": attr.string(default = "Hello, World!"),
    },
    outputs = {"out": "%{name}.txt"},
)

def print_hello():
    """Prints a hello world message"""
    print("Hello from Bazel!")