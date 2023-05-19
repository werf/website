# Prepare system

```
# /etc/hosts

127.0.0.1       k3d-registry.sample-app.test
127.0.0.1       sample-app.test
127.0.0.1       sample-app-2.test
```

# Run scenarios

## Successful deployment: werf vs helm

Prepare:

```
./scripts/prepare_successful_deploy.sh
```

Run scenario:

```
./scripts/run_successful_deploy.sh
```

## Deployment with failure: werf vs helm

Prepare:

```
./scripts/prepare_deploy_with_failure.sh
```

Run scenario:

```
./scripts/prepare_deploy_with_failure.sh
```

## Build, deploy, test, cleanup and distribute demos

Required terminal window size:

```
height: 644px;
width: 760px;
```

Required terminal font size: 9pt.
