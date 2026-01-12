# test-check

The test-check is a minimal Kuberhealthy check that reports a configured
success or failure after a configurable delay. It is useful for validating
that Kuberhealthy reporting works end-to-end.

## Configuration

The check reads the following environment variables:

- `REPORT_FAILURE` (optional, default: `false`): When set to `true`, the check
  reports a failure to Kuberhealthy.
- `REPORT_DELAY` (optional, default: `5s`): Duration to wait before reporting.

Kuberhealthy also injects required reporting variables such as
`KH_REPORTING_URL` and `KH_CHECK_RUN_DEADLINE`.

## Run the check

1. Update the image tag in `healthcheck.yaml`.
2. Apply the example HealthCheck to your cluster:

```sh
kubectl apply -f healthcheck.yaml
```

## Build locally

```sh
just build
```

## Test locally

```sh
just test
```
