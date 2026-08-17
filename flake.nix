{
  inputs = {
    nixpkgs.url = "tarball+https://git.tatikoma.dev/corpix/nixpkgs/archive/corpix.tar.gz";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        go = pkgs.go_1_27.overrideAttrs (_: rec {
          version = "1.27rc3";

          src = pkgs.fetchurl {
            url = "https://go.dev/dl/go${version}.src.tar.gz";
            hash = "sha256-6eIO3RcgCV+RCWluljpmBp0/bUjQDrk4jIiM4GYx31w=";
          };
        });
      in
      {
        packages.default = go;

        devShells.default = pkgs.mkShell {
          packages = [
            go
          ];

          # Never silently download another Go toolchain.
          GOTOOLCHAIN = "local";
        };
      }
    );
}
