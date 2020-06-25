# Creamy Artifacts

<a href="https://hub.docker.com/r/albinodrought/creamy-artifacts">
<img alt="albinodrought/creamy-artifacts Docker Pulls" src="https://img.shields.io/docker/pulls/albinodrought/creamy-artifacts">
</a>
<a href="https://github.com/AlbinoDrought/creamy-artifacts/blob/master/LICENSE"><img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/creamy-artifacts"></a>

Store and merge build artifacts.

## Usage

Right now there are no configuration options.

See [example/run.sh](./example/run.sh) for HTTP usage.

### With Docker

```sh
docker run --rm -it \
    -v $(pwd)/foo/bar:/data \
    albinodrought/creamy-artifacts
```

### Without Docker

```sh
./creamy-artifacts
```

## Building

### With Docker

```sh
make image
```

### Without Docker

```sh
make build
```
