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

All build steps are listed in the build file inside `tasks:` array.
The steps are run in the same order they listed in the build file.

### Step configuration

Each step in run inside a Docker container.
Enter the Docker image in `image` field.
The `script` field holds an array of all the commands to be run inside the container.
The working directory will be available in the container under `/source/` directory.
By default the working directory is the directory where the `comake` command is run in.
It can be changed using `--workdir` option.
Please note that, for now, if you need to use scripting features you need to use `sh -c "<script>"` (or similar for diffent shells).
Better shell support will be added in the future.