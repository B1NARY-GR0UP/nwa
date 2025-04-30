{
  description = "A simple hello world flake";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.${system}.hello = pkgs.writeShellScriptBin "hello" ''
        echo "Hello, World!"
      '';

      defaultPackage.${system} = self.packages.${system}.hello;
    };
}