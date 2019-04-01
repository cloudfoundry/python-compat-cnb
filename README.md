# Python-compat Cloud Native Buildpack

This buildpack handles some of the V2B python specific logic and is designed to allow for parity between earlier buildpacks and the Python language CNBs as a whole.

---
To package this buildpack for consumption:
```
$ ./scripts/package.sh
```
This builds the buildpack's Go source using GOOS=linux by default. You can supply another value as the first argument to package.sh.
