# Skaffold Environment Templating

Skaffold supports environment templating to make configuration files more flexible and reusable. This allows you to use environment variables and built-in functions to dynamically set values in your Skaffold configuration files.

## Summary Table

| Feature                | Syntax/Function           | Description                                                                 |
|------------------------|--------------------------|-----------------------------------------------------------------------------|
| Environment Variable   | `${VAR}`                 | Inserts value of `VAR` from environment                                     |
| Default Value          | `${VAR:default}`         | Uses `default` if `VAR` is not set                                          |
| Get Env Function       | `${env:VAR}`             | Fetches value of `VAR` from environment                                     |
| Required Env Function  | `${requiredEnv:VAR}`     | Fails if `VAR` is not set                                                   |
| Expand Env Function    | `${expandenv:STRING}`    | Expands all environment variables in `STRING`                               |

## Key Features

- **Environment Variables**: You can reference environment variables in your Skaffold YAML files using the `${VAR}` or `${VAR:default}` syntax.
- **Built-in Functions**: Skaffold provides functions like `env`, `requiredEnv`, and `expandenv` to fetch and manipulate environment variables.
- **Default Values**: You can specify default values for variables if they are not set in the environment.
- **Required Variables**: Use `requiredEnv` to enforce that a variable must be set, otherwise Skaffold will fail with an error.

## Examples

```yaml
apiVersion: skaffold/v4beta6
kind: Config
metadata:
  name: my-app
build:
  artifacts:
    - image: my-image
      context: .
      docker:
        dockerfile: Dockerfile
deploy:
  kubectl:
    manifests:
      - k8s/*.yaml
profiles:
  - name: dev
    activation:
      - env: DEV
    build:
      artifacts:
        - image: my-image
          context: .
          docker:
            dockerfile: Dockerfile
    deploy:
      kubectl:
        manifests:
          - k8s/dev/*.yaml
```

You can use environment variables like this:

```yaml
build:
  artifacts:
    - image: my-image:${TAG}
```

Or with a default value:

```yaml
build:
  artifacts:
    - image: my-image:${TAG:latest}
```

## Useful Functions

- `${env:VAR}`: Gets the value of `VAR` from the environment.
- `${requiredEnv:VAR}`: Requires `VAR` to be set, otherwise fails.
- `${expandenv:STRING}`: Expands all environment variables in `STRING`.

## References

For more details, see the [Skaffold Environment Templating documentation](https://skaffold.dev/docs/environment/templating/).
