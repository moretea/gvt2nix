{ pkgs ? (import <nixpkgs> {}) }: 
with pkgs;
with (callPackage ../../. {});
buildGvtPackage {
  name = "catalogue";

  src = fetchFromGitHub {
    owner = "microservices-demo";
    repo = "catalogue";
    rev = "93452255c62cac6e44f5d86376d7cb3c5409e162";
    sha256 = "1f0jv6gvax6aijdl4b8z0v1g5bvb0pj0i7l51cpi36jgnhrklfng";
  };

  depsSha256 = "003ibzn63hccg0fhdvf8n0gzx281i6npbhl7c0y023yypj5w5s24";
  goPackagePath = "github.com/microservices-demo/catalogue";
}

