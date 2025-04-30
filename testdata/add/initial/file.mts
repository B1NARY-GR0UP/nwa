export function sayHello(name: string): string {
    return `Hello, ${name}!`;
}

export interface Greeting {
    message: string;
    recipient: string;
}

export const defaultGreeting: Greeting = {
    message: "Hello",
    recipient: "World"
};