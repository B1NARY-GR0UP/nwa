function sayHello(name: string): string {
    return `Hello, ${name}!`;
}

const greeting: string = "World";

module.exports = {
    sayHello,
    greeting
};