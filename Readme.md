# kpg

WARNING: proof of concept, not for general consumption

## Test locally

```sh
# install
$ go get github.com/ecordell/kpg

# start a registry
$ docker run -it --rm -p 5000:5000 registry

# push manifests
$ kpg push ./example localhost:5000/oras:test
Pushing to localhost:5000/oras:test...
Pushed  with digest sha256:58e5db1255a64fa55754a04ec3b1ef0376626266b0b169cb4bec4f8aa92220c8

# pull manifests
$ kpg pull localhost:5000/oras:test -o ./out
Pulling from localhost:5000/oras:test and saving...
Pulled from localhost:5000/oras:test with digest sha256:6870ed324c56aabbd016ea1252275761a1a9610f60d4ac3c2df933e7b532a560

# by digest
$ kpg pull localhost:5000/oras@sha256:6870ed324c56aabbd016ea1252275761a1a9610f60d4ac3c2df933e7b532a560 -o ./out
Pulling from localhost:5000/oras@sha256:6870ed324c56aabbd016ea1252275761a1a9610f60d4ac3c2df933e7b532a560 and saving...
Pulled from localhost:5000/oras@sha256:6870ed324c56aabbd016ea1252275761a1a9610f60d4ac3c2df933e7b532a560 with digest sha256:6870ed324c56aabbd016ea1252275761a1a9610f60d4ac3c2df933e7b532a560

$ tree out
  out
  ├── configMap.yaml
  ├── deployment.yaml
  ├── kustomization.yaml
  └── service.yaml
```

## MediaTypes

Two media types are defined - one for any kube yaml, and one for the `kustomization.yaml` file. JSON will be supported
in the future.

See them in action via logs:

```sh
$ kpg pull localhost:5000/oras:test ./out                                                                            130 ↵
Pulling from localhost:5000/oras:test and saving...
WARN[0000] unknown type: application/vnd.oci.image.config.v1+json
WARN[0000] reference for unknown type: application/vnd.k8s.manifest.v1+yaml  digest="sha256:8793d287579a9ed1590c41c600c60fa634d205523b56860425f38102a00e12fb" mediatype=application/vnd.k8s.manifest.v1+yaml size=183
WARN[0000] reference for unknown type: application/vnd.k8s.kustomization.manifest.v1+yaml  digest="sha256:d8b8f45bdf793b692f3d1511bfda1e1170531ef4eb8268eea8b084d8ce38c997" mediatype=application/vnd.k8s.kustomization.manifest.v1+yaml size=89
WARN[0000] reference for unknown type: application/vnd.k8s.manifest.v1+yaml  digest="sha256:eef83564b70c404ada0c43c28f155f07df025cbe889e4ead32c612c7af2a6910" mediatype=application/vnd.k8s.manifest.v1+yaml size=697
WARN[0000] reference for unknown type: application/vnd.k8s.manifest.v1+yaml  digest="sha256:78d9149b6e67fdcdf395069f8497cab474c95033a5733beddda7858e6b9dbd24" mediatype=application/vnd.k8s.manifest.v1+yaml size=117
WARN[0000] encountered unknown type application/vnd.k8s.kustomization.manifest.v1+yaml; children may not be fetched
WARN[0000] encountered unknown type application/vnd.k8s.manifest.v1+yaml; children may not be fetched
WARN[0000] encountered unknown type application/vnd.k8s.manifest.v1+yaml; children may not be fetched
WARN[0000] encountered unknown type application/vnd.k8s.manifest.v1+yaml; children may not be fetched
Pulled from localhost:5000/oras:test with digest sha256:6870ed324c56aabbd016ea1252275761a1a9610f60d4ac3c2df933e7b532a560
```
