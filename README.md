# Todoer

Todoer is a service responsible for allowing clients to create TODO lists.

This service can also be served via [gRPC](https://grpc.io/).

But the default option is to use [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) and [JSON](https://www.json.org/json-en.html) format.
For wich the reference documentation can be found [here](api.md).

## Development

### Dependencies

To run the automation provided here you need to have installed:

* [Make](https://www.gnu.org/software/make/)
* [Docker](https://docs.docker.com/get-docker/)

It is recommended to just use the provided automation through make,
it will help you achieve consistent results across different hosts
and even different operational systems.

If you fancy building things (or running tests) with no extra layers
you can install the latest [Go](https://golang.org/doc/install) and run
Go commands like:

```sh
go test -race ./...
```

### Running tests

Just run:

```
make test
```

To check locally the coverage from all the tests run:

```
make coverage
```

And it should open the coverage analysis in your browser.


### Releasing

To create an image ready for production deployment just run:

```
make image
```

And one will be created tagged with the git short revision of the
code used to build it, you can also specify an explicit version
if you want:

```
make image version=0.0.1
```

### Running Locally

If you want to explore a locally running version of the service just run:

```
make run
```

And the service will be available at port 8080.

You can also specify the options:

```
make run opts='-grpc'
```

### Compiling a binary

If you want to compile a binary for your host system and you have the `go` tool
installed, you can compile your own binary by running:

```
make build
```

It will generate binary at `cmd/todoer/todoer` with the following options:

```
Usage of ./cmd/todoer/todoer:
  -grpc
      run todoer service with grpc server
  -port int
      port where the service will be listening to (default 8080)
```

### Generating Protobuf and gRPC code

You can change the `pb/todoer.proto` file and run:

```
make generate
```

It will generate `pb/todoer.pb.go` and `pb/todoer_grpc.pb.go` with the structs
and functions to be used by the api server part of the code.

### UI to check gRPC functions

If you have [grpcui]() installed, you can test the exposed gRPC functions via
user interface in your browser.

First you will need to be running todoer service:

```
make run opts='-grpc'
```

Then in another shell session you can run:

```
make grpcui
```

## Deploy

This application has a deploy strategy to a [Kubernetes](https://kubernetes.io/) cluster.
It uses [Kustomize](https://kustomize.io/) for templating.
You can write your own `kustomization.yaml` to a directory `deploy` with different variables like this:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- github.com/vitorarins/todoer/deploy

images:
- name: vitorarins/todoer
  newName: your-image-name
  newTag: your-image-tag
```

And apply to your cluster by running: `kubectl apply -k deploy`.

If you are using this repository you can change `kustomization.yaml` on the `deploy` directory and run:

```
make deploy version=0.0.1
```

And the service will be available through a `LoadBalancer` type of Kubernetes Service.
