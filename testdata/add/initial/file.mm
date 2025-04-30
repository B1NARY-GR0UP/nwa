#import <Foundation/Foundation.h>
#include <iostream>

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSLog(@"Hello, World from Objective-C!");
        std::cout << "Hello, World from C++!" << std::endl;
    }
    return 0;
}