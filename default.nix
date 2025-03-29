{ buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "hll-arty-tui";
  version = "0.0.1";

  src = fetchFromGitHub {
    owner = "computerdane";
    repo = "hll-arty-tui";
    rev = "v${version}";
    hash = "sha256-9bl+5Bk/xMSuf1iX5Vy3tZqCrATnpwIkMQeNERRJ+/U=";
  };

  vendorHash = "sha256-3NehgSR6y64/6VlvS7CcdAcekHTGXOHgWW0Zn3uX8hc=";
}
