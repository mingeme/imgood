{
  description = "Imgood";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "imgood";
            version = "0.1.0";
            src = ./.;
            nativeBuildInputs = [ pkgs.pkg-config ];
            buildInputs = [ pkgs.vips ];

            vendorHash = "sha256-lRm3D9xAIrr6BnY5wRzJaMlzEnNfxr5LFWVGbWmprwI=";
          };
        };

        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/imgood";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            vips
            pkg-config
          ];

          shellHook = ''
            echo "Imgood development environment"
            echo "Go version: $(go version)"
            echo "Vips version: $(pkg-config --modversion vips)"
            echo ""
            echo "Ready to build and run the imgood!"
          '';
        };
      }
    );
}
