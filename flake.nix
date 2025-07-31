{
  description = "Minimal development environment with mise";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        # Minimal tools - mise will handle language runtimes and dev tools
        minimalTools = with pkgs; [
          # Core system tools
          curl
          wget
          git
          jq

          # mise for tool version management
          mise

          # Docker for containerization
          docker
          docker-compose

          # Basic database clients (commonly needed)
          postgresql_17
          redis
        ];

      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = minimalTools;

          shellHook = ''
            echo "🚀 Minimal Development Environment"
            echo ""
            echo "Available tools:"
            echo "  • mise (managing language runtimes and dev tools)"
            echo "  • Docker & Docker Compose"
            echo "  • PostgreSQL & Redis clients"
            echo "  • Basic utilities (curl, wget, git, jq)"
            echo ""
            echo "mise will manage:"
            echo "  • Go, Node.js, pnpm"
            echo "  • Development tools (air, golangci-lint, terraform, etc.)"
            echo ""
            echo "Run 'mise install' to install all configured tools"
            echo ""

            # Ensure mise is activated
            eval "$(mise activate bash)"
          '';

          # Environment variables
          NIX_SHELL_PRESERVE_PROMPT = 1;
        };

        # Formatter for the flake
        formatter = pkgs.nixpkgs-fmt;
      });
}
