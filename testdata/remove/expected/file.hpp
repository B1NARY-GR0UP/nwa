#pragma once

#include <string>
#include <vector>

namespace hello {
    class HelloWorld {
    public:
        HelloWorld();
        std::string getMessage() const;
    private:
        std::vector<std::string> message;
    };
}