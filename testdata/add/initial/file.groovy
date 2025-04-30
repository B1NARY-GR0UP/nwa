class HelloWorld {
    static void main(String[] args) {
        println "Hello, World!"

        def greeter = new Greeting()
        greeter.sayHello("Groovy")
    }
}

class Greeting {
    def sayHello(name) {
        println "Hello, ${name}!"
    }
}