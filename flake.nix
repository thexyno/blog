{
  description = "A simple Go package";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-22.11"; # Tailwind isn't in 21.11

  outputs = { self, nixpkgs }:
    let

      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    {
      nixosModule = { config, options, lib, pkgs, ... }:
        let
          cfg = config.services.xynoblog;
          xb = self.packages.${pkgs.system}.xynoblog;
        in
        with lib;
        {
          options.services.xynoblog = {
            enable = mkOption {
              type = types.bool;
              default = false;
              description = "wether to enable the xynoblog blog engine";
            };
            listen = mkOption {
              type = types.str;
              default = ":8392";
              description = "the domain/post xynoblog listens on";
            };
            stateDirectory = mkOption {
              type = types.str;
              default = "xynoblog";
              description = "dir to store the sqlite3 database (relative to /var/lib)";
            };

          };
          config = mkIf cfg.enable {
            users.users.xynoblog = {
              group = "xynoblog";
              isSystemUser = true;
            };
            users.groups.xynoblog = { };

            environment.systemPackages = [
              (pkgs.writeScriptBin "xynoblog" ''
                #!/bin/sh
                exec ${xb}/bin/xynoblog --db /var/lib/${cfg.stateDirectory}/blog.db --mediadir /var/lib/${cfg.stateDirectory}/media "$@"
              '')
            ];

            systemd.services.xynoblog = {
              description = "xynoblog blog engine";
              after = [ "network.target" ];
              wantedBy = [ "multi-user.target" ];
              serviceConfig = {
                User = "xynoblog";
                Group = "xynoblog";
                PrivateTmp = "true";
                PrivateDevices = "true";
                ProtectHome = "true";
                ProtectSystem = "strict";
                AmbientCapabilities = "CAP_NET_BIND_SERVICE";
                ExecStart = "${xb}/bin/xynoblog serve --listen \"${cfg.listen}\" --db /var/lib/${cfg.stateDirectory}/blog.db --mediadir /var/lib/${cfg.stateDirectory}/media";
                StateDirectory = cfg.stateDirectory;
              };
            };
          };
        };
      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          xynoblog_tmpl =
            pkgs.mkYarnPackage rec {
              pname = "xynoblog_tmpl";
              inherit version;
              offlineCache = pkgs.fetchYarnDeps {
                yarnLock = src + "/yarn.lock";
                sha256 = "sha256-mSxDAI1PQ3muGnNbkqha4nV5S5htk4gaO6gbh/Z3Zfk=";
              };
              src = ./.;
              packageJSON = ./package.json;
              distPhase = "true";
              buildPhase = ''
                export HOME=$(mktemp -d)
                echo $node_modules
                yarn --offline build
              '';
            };
          xynoblog =
            pkgs.buildGo118Module rec {
              pname = "xynoblog";
              inherit version;
              # In 'nix develop', we don't need a copy of the source tree
              # in the Nix store.
              src = ./.;
              nativeBuildInputs = [ pkgs.installShellFiles pkgs.makeWrapper ];
              buildInputs = [ pkgs.libwebp ];

              preBuild = ''
                cp -r ${self.packages.${pkgs.system}.xynoblog_tmpl}/* .
                go generate ./...
              '';

              preFixup = ''
                wrapProgram $out/bin/xynoblog \
                  --prefix XYNOBLOG_STATICDIR : "$out/share/xynoblog" \
                  --prefix GIN_MODE : "release" \
                  --prefix XYNOBLOG_RELEASEMODE : "true"
              '';
              postInstall = ''
                mkdir -p $out/share/xynoblog
                cp -r ./data $out/share/xynoblog/data
                installShellCompletion --cmd ${pname} \
                  --bash <($out/bin/${pname} completion bash) \
                  --fish <($out/bin/${pname} completion fish) \
                  --zsh  <($out/bin/${pname} completion zsh)
              '';

              # This hash locks the dependencies of this package. It is
              # necessary because of how Go requires network access to resolve
              # VCS.  See https://www.tweag.io/blog/2021-03-04-gomod2nix/ for
              # details. Normally one can build with a fake sha256 and rely on native Go
              # mechanisms to tell you what the hash should be or determine what
              # it should be "out-of-band" with other tooling (eg. gomod2nix).
              # To begin with it is recommended to set this, but one must
              # remeber to bump this hash when your dependencies change.
              #vendorSha256 = pkgs.lib.fakeSha256;

              vendorSha256 = "sha256-EgM7902IswosjKuyLy4aJEuVwVNEveUXexjbWoljIcc=";
            };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.xynoblog);
      devShell = forAllSystems (system:
        let pkgs = nixpkgsFor.${system}; in
        (pkgs.mkShell {
          XYNOBLOG_FONTDIR = "${pkgs.jetbrains-mono}/share/fonts/truetype";
          buildInputs = [ pkgs.nixpkgs-fmt pkgs.gopls pkgs.go_1_18 pkgs.lefthook pkgs.libwebp pkgs.yarn2nix pkgs.yarn ];
        }));
    };
}
