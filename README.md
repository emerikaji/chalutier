# Chalutier

Small utility to get the version of a Docker image from its Docker Hub tags.

## Dependencies

chalutier can run without Docker in image digest mode, as it will directly compare against the registry.

To use its container ID mode, Docker needs to be installed and running.

## Usage

```
chalutier (<container id> | -d <image digest>) [-r <registry>]
```

For example
```bash
$ chalutier -d alpine@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0
```
Outputs
```bash
3.20.1
```

If the container uses a locally built image, or if the digest cannot be found in the tags, `chalutier` will simply output
```
null
```

The default registry that will be used if the argument isn't provided is `https://hub.docker.com`.

## Installation

### Build from source (All platforms)

After cloning the project, simply run
```bash
CGO_ENABLED=0 go build
```

The `chalutier` utility will be built inside the current folder.

This project was made with Go version 1.22.5, so no guarantee can be provided that it will build on earlier versions. However, any Go version recent enough to support [Moby](https://pkg.go.dev/github.com/docker/docker) should work.