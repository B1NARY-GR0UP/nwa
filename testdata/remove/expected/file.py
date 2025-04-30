def say_hello(name="World"):
    """Print a hello message"""
    return f"Hello, {name}!"

class Greeter:
    def __init__(self, greeting="Hello"):
        self.greeting = greeting

    def greet(self, name):
        return f"{self.greeting}, {name}!"

if __name__ == "__main__":
    print(say_hello())