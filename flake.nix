{
  description = "Interactive tmux session manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      eachSystem = flake-utils.lib.eachDefaultSystem;

      # Bump this when tagging a new release
      version = "1.0.4";

      mkPkg = pkgs: (pkgs.buildGoModule {
        pname = "mux-session";
        inherit version;
        src = ./.;
        vendorHash = "sha256-7it/kFzW4DG3Cmvevipe5oxf9KHTnR4KyB8C6aE903Q=";
        ldflags = [ "-s" "-w" ];
        subPackages = [ "." ];
        meta = with pkgs.lib; {
          description = "Interactive tmux session manager";
          homepage = "https://github.com/niedch/mux-session";
          mainProgram = "mux-session";
          platforms = platforms.linux ++ platforms.darwin;
        };
      }).overrideAttrs (old: {
        env = (old.env or {}) // { CGO_ENABLED = "0"; };
      });
    in
    eachSystem (system: let pkgs = import nixpkgs { inherit system; }; in {
      packages.default = mkPkg pkgs;

      devShells.default = pkgs.mkShell {
        packages = with pkgs; [ go gopls golangci-lint tmux fzf ];
      };
    }) // {
      overlays.default = final: prev: {
        mux-session = mkPkg final;
      };
    };
}
