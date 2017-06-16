# gvt2nix

With `gvt2nix`, you can build Go programs that use gvt to manage their dependencies.


## Impression

```nix
buildGvtPackage {
  name = "catalogue";
  src = fetchFromGitHub { ... };
  depsSha256 = "003ibzn63hccg0fhdvf8n0gzx381i6npbhl7c0y023yypj5w5s24";
  goPackagePath = "github.com/myorg/myprog";
}
```

## Examples
See the `examples/` directory
