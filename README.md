# Chaos Proxy

`chaosproxy` is a configurable proxy written in golang. It is meant to be used in game day testing and drills that test how products handle outages of dependencies and external services that the service uses.

It can be used to simulate random outages in the same vein as Netflix's Chaos Monkey project. This tool differs in that it is primarily concerned with introducing chaos by sitting in between your service and its external APIs. It's useful for testing the resiliency of microservice architectures.

## Installation and use

Run directly:
```bash
$ dep ensure
$ go run main.go
```

Using docker:
```bash
$ docker build -t chaosproxy .
$ docker run -v ${PWD}/test_config.yaml:/chaosproxy.yaml -t chaosproxy
```

The server listens on `:8080`, so if you set your `HTTP_PROXY` environment variable or equivalent setting in your service, all requests will go through chaosproxy.

## Configuration

`chaosproxy` is configured using "scripts" in YAML. Scripts consist of a regex pattern to match incoming requests and a list of actions, which are chained in the order that they appear in the file for the matching request.

The list of actions is where you introduce chaos. Some actions delay responses for random amounts of time, some of them return errors outright.

### Script example

This script sets up a proxy that intercepts a service's calls to `google.com` and introduces a random delay of 1-4 seconds:

```yaml
- pattern: google.com
  actions:
    - name: randsleep
      params:
        from: 1
        to: 4
```

### Actions

See [currently supported actions](docs/actions.md) for a list of what actions are available right now to use.

## Contributing

You can fork, branch, and modify this codebase to your liking. If you'd like to contribute, see the [contribution guidelines](CONTRIBUTION.md).
