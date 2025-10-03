{
  description = "A simple tetris game in the terminal";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "pterm-tetris";
          version = "0.0.1";

          src = ./.;

          vendorHash = "sha256-yWSdZLfdV7rbCWZgnpm/Y8ZwzNYIC/RreoY3vTTtgvE=";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
          ];
        };
      });
}
