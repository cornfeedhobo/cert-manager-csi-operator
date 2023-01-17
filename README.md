# cert-manager-csi-operator

Admission Webhook example, using [`operator-sdk`](https://github.com/operator-framework/operator-sdk).

This is to demonstrate a working example of using the `operator-sdk` to manage
the lifecycle of an operator that is implementing a core type,
or maybe a type managed by another repository.

## Usage

Each commit logs the purpose and relevant commands.

```shell
git log \
  --date=format:'%a %Y-%m-%d %k:%M' \
  --format='%n%ad %an <%ae> %h%d%n%n %s%n%w(0,2,2)%+b' \
  --compact-summary
```

## To Do

- [x] Initialize with `operator-sdk`

- [x] Scaffold a core type

- [x] Make sure kustomize can render manifests

- [x] Remove all the temporary API code, and get tests working to validate

- [x] Add some more core types and add some corresponding tests

- [x] Add helm chart, created from the kustomize output

- [ ] Add auth support for the metrics endpoint

- [ ] Add support for leader election

- [ ] Add instructions for testing without cert-manager
