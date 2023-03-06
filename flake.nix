{
  description = "Development environment";

  outputs = { self, nixpkgs }: {
    devShells.aarch64-darwin.default = with nixpkgs.legacyPackages.aarch64-darwin; mkShell {
      nativeBuildInputs = [
        go
        gopls
        terraform
        golangci-lint
      ];
    };

  };
}
