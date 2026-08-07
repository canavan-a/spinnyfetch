{
  description = "spinnyfetch - a NixOS system-info fetcher with a spinning Nix logo";

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
          pname = "spinnyfetch";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-CJ33GAE7d87w7Ld1hTcjuz34gpqdagEOWx7ziCCAywQ=";

          meta = with pkgs.lib; {
            description = "NixOS system-info fetcher with a spinning Nix logo";
            homepage = "https://github.com/canavan-a/spinnyfetch";
            license = licenses.mit;
            mainProgram = "spinnyfetch";
            platforms = platforms.linux;
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [ pkgs.go pkgs.gopls ];
        };
      });
}
