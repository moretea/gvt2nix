{
  cacert,
  stdenv,
  buildGoPackage,
  fetchFromGitHub,
  git
}: rec {
  gvt = buildGoPackage rec {
    name = "gvt-${version}";
    rev = "50d83ea21cb0405e81efd284951e111b3a68d701";
    version = builtins.substring 0 7 rev;

    src = fetchFromGitHub {
      owner = "FiloSottile";
      repo = "gvt";
      inherit rev;
      sha256 = "0z9a7yv9wpj3mby44nkpdrk1s2v3slir96ap60i0lw08hr0krd6y";
    };

    goPackagePath = "github.com/FiloSottile/gvt";

    goDeps = null ;
  };

  gvtRestored = {name, manifest, sha256}: stdenv.mkDerivation {
    inherit name;

    buildInputs = [ gvt git];

    outputHashAlgo = "sha256";
    outputHashMode = "recursive";
    outputHash = sha256;

    GIT_SSL_CAINFO = "${cacert}/etc/ssl/certs/ca-bundle.crt";

    buildCommand = ''
      mkdir -p src/pkg/vendor
      cp ${manifest} src/pkg/vendor/manifest
      GOPATH=`pwd`
      cd src/pkg
      gvt restore
      mkdir -p $out/vendor
      cp -r vendor/* $out/vendor/
    '';
  };

  gvtSource = name: src: gvtVendor: stdenv.mkDerivation {
    inherit name;

    buildCommand = ''
      mkdir $out
      cp -vr ${src}/* $out/
      chmod -R +w $out/vendor
      cp -r ${gvtVendor}/vendor/* $out/vendor
      chmod -w $out/vendor
    '';
  };

  buildGvtPackage = { name, src, depsSha256, goPackagePath }: 
  let
    goSrc = src;
  in buildGoPackage rec {
    inherit name;
    inherit goPackagePath;

    gvtRestoredDeps = gvtRestored {
      name = "${name}-gvt-deps";
      manifest = goSrc + /vendor/manifest;
      sha256 = depsSha256;
    };

    src = gvtSource "${name}-src" goSrc gvtRestoredDeps;

    deps = null;
  };
}
