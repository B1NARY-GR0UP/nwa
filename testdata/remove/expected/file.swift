import Foundation

class Greeter {
    let greeting: String

    init(greeting: String) {
        self.greeting = greeting
    }

    func greet() {
        print(greeting)
    }
}

let greeter = Greeter(greeting: "Hello, World!")
greeter.greet()