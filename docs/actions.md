# Chaos Proxy Actions

Below is a list of the available actions you can use to configure Chaos Proxy. See the [contrib docs for proxy actions](../CONTRIBUTION.md#actions) for how to add more!

## Sleep

Sleep for _n_ seconds before continuing to process the request.

* Action name: `sleep`
* Params: `seconds (int)`

Example:
```yaml
name: sleep
params:
  - seconds: 5
```

## Random Sleep

Sleep for a random interval between _n_ and _m_ seconds before continuing to process the request.

* Action name: `sleeprand`
* Params: `from (int); to (int)`

Example:
```yaml
name: sleeprand
params:
  - from: 1
  - to: 10
```

## Error

Return an HTTP error with a custom message.

* Action name: `httperror`
* Params: `code (int); status (str)`

Example:
```yaml
name: httperror
params:
  - code: 500
  - status: Internal Server Error
```
