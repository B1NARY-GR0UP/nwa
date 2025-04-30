#ifndef HELLO_WORLD_HH
#define HELLO_WORLD_HH

#include <string>

namespace hello {
    class Greeter {
    public:
        Greeter();
        std::string greet();
    };
}

#endif // HELLO_WORLD_HH