# Todoer

Todoer is a service responsible for allowing clients to create TODO lists.

The reference documentation for the API can be found [here](api.md).

## Development

### Dependencies

To run the automation provided here you need to have installed:

* [Make](https://www.gnu.org/software/make/)
* [Docker](https://docs.docker.com/get-docker/)

It is recommended to just use the provided automation through make,
it will help you achieve consistent results across different hosts
and even different operational systems (it is also used on the CI).

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
make image version=1.12.0
```

### Running Locally

If you want to explore a locally running version of the service just run:

```
make run
```

And the service will be available at port 8080.

### Deploy

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
make deploy
```

And the service will be available through a `LoadBalancer` type of Kubernetes Service.
