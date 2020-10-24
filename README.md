# comake - Containerized Build Tool

Comake enables containerization for builds. Containerization for runtime is a good idea, it should be same for buildtime too.

## Example

Given `hello-build.yml` as follows:

```yaml
steps:

- name: greeting
  image: docker.io/library/alpine
  script:
  - echo "starting to greet everyone"
  - sh -c 'mkdir -p build'
  - sh -c 'echo "hello world" > build/hello'
```


```bash
comake -f hello-build.yml run
```

## Buildfile documentation

### Build steps

All build steps are listed in the build file under `steps:` array.
The steps are ran one after another in the same order they were declared.

### Step configuration

Each step is run inside a Docker container.
The image for the Docker container should be provided under `image` field.
The `script` field holds an array of all the commands to be run inside the container.

The working directory will be available in the container under `/source/` directory.
By default the working directory is the directory where the `comake` command is run in.
It can be changed using `--workdir` option.

Please note that, as of now, if you need to use scripting features you need to use `sh -c "<script>"` (or similar for diffent shells).
Better shell support will be added in the future.