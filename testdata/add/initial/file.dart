void main() {
  print('Hello, World!');

  var greeting = Greeting('Hello');
  greeting.printMessage('World');
}

class Greeting {
  String message;

  Greeting(this.message);

  void printMessage(String name) {
    print('$message, $name!');
  }
}