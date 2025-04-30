module hello_world;
  string message = "Hello, World!";

  initial begin
    $display(message);
    print_greeting();
    #10 $finish;
  end

  function void print_greeting();
    $display("Greetings from SystemVerilog!");
  endfunction
endmodule