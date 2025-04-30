function greet(name: string): string {
    return `Hello, ${name}!`;
}

class Greeter {
    greeting: string;

    constructor(message: string) {
        this.greeting = message;
    }

    greet(): string {
        return `Hello, ${this.greeting}!`;
    }
}

console.log(greet("World"));