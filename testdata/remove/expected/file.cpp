#include <iostream>
#include <string>

class Greeter {
private:
    std::string message;
public:
    Greeter(const std::string& msg) : message(msg) {}
    void greet() {
        std::cout << message << std::endl;
    }
};

int main() {
    Greeter greeter("Hello, World!");
    greeter.greet();
    return 0;
}