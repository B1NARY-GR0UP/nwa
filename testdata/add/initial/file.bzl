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