# envkey-source

Integrate [EnvKey](https://www.envkey.com) with any language, either in development or on a server, by making your configuration available through the shell as environment variables.

# v2

Now that [EnvKey v2](https://v2.envkey.com) has been released, you can find version 2 of envkey-source in [a subdirectory of the EnvKey v2 monorepo](https://github.com/envkey/envkey/tree/main/public/sdks/envkey-source). Using v2 requires an EnvKey v2 organization (it won't work with ENVKEYs generated in a v1 org).

[Here's a guide on migrating from v1 to v2.](https://docs-v2.envkey.com/docs/migrating-from-v1)

## Installation

envkey-source compiles into a simple static binary with no dependencies, which makes installation a simple matter of fetching the right binary for your platform and putting it in your `PATH`. An `install.sh` script is available to simplify this, as well as a [homebrew tap](https://github.com/envkey/homebrew-envkey).

**Install via bash:**

```bash
curl -s https://raw.githubusercontent.com/envkey/envkey-source/master/install.sh | bash
```

***Note:** the install.sh script writes, then deletes a couple temporary files to the current directory during installation, so make sure you have write permissions for whatever directory you run this command in. In locked down environments, you may want to run it in `$HOME` to be safe.*

***Another Note:** the install.sh script downloads the appropriate envkey-source binary from Github Releases by default, but Github has a fairly low rate limit for unauthenticated requests. For added redundancy, it will fail over to envkey-source-releases.s3.amazonaws.com (an S3 bucket controlled by EnvKey) and download the binary from there if the request to Github Releases fails.*

**Install manually:**

Find the [release](https://github.com/envkey/envkey-source/releases) for your platform and architecture, and stick the appropriate binary somewhere in your `PATH` (or wherever you like really).

## Usage

First, generate an `ENVKEY` in the [EnvKey App](https://github.com/envkey/envkey-app).

Then with a `.env` file in the current directory that includes `ENVKEY=...` (in development) / an `ENVKEY` environment variable set (on a server):

```bash
eval $(envkey-source)
```

Now you can access your app's environment variables in this shell, or in any process (in any language) launched from this shell.

You can also pass an `ENVKEY` directly. This isn't recommended for a real workflow, but can be useful for trying things out.

```bash
eval $(envkey-source ENVKEY)
```

### Multi-Line Values

If your EnvKey config includes multi-line values, you need to load it with slightly different syntax to preserve formatting. Instead of:

```bash
eval $(envkey-source)
echo $SOME_VAR
```

Use this (note the additional **double quotes**):

```bash
eval "$(envkey-source)"
echo "$SOME_VAR"
```

### Flags

```text
    --cache                cache encrypted config as a local backup (default is true when .env file exists, false otherwise)
    --no-cache             do NOT cache encrypted config as a local backup even when .env file exists
    --cache-dir string     cache directory (default is $HOME/.envkey/cache)
    --dot-env-compatible   change output to .env format
    --env-file string      ENVKEY-containing env file name (default ".env")
    --pam-compatible       change output format to be compatible with /etc/environment on Linux
-f, --force                overwrite existing environment variables and/or other entries in .env file
-h, --help                 help for envkey-source
-v, --version              prints the version
    --verbose              print verbose output (default is false)
    --timeout float        timeout in seconds for http requests (default 10)
    --retries uint8        number of times to retry requests on failure (default 3)
    --retryBackoff float   retry backoff factor: {retryBackoff} * (2 ^ {retries - 1}) (default 1)
```

### Errors

If you get an error, envkey-source will echo the error string to stdout and return false instead of setting environment variables. For example:

```bash
$ eval $(envkey-source notvalidenvkey) && ./env-dependent-script.sh
error: ENVKEY invalid
```

### Security - Preventing Shell Injection

Whenever you use `eval`, you need to worry about shell injection. We did the worrying for you--envkey-source wraps all EnvKey variables in single quotes and safely escapes any single quotes the variables might contain. This removes any potential for shell injection.

### Overriding Vars

By default, envkey-source will not overwrite existing environment variables or additional variables set in a `.env` file. This can be convenient for customizing environments that otherwise share the same configuration. But if you do want EnvKey vars to take precedence, use the `--force` / `-f` flag. You can also use [sub-environments](https://blog.envkey.com/development-staging-production-and-beyond-85f26f65edd6) in the EnvKey App for this purpose.

### Working Offline

envkey-source caches your encrypted config in development so that you can still use it while offline. Your config will still be available (though possibly not up-to-date) the next time you lose your internet connection. If you do have a connection available, envkey-source will always load the latest config.

By default, caching is enabled when a `.env` file is present in the directory, and disabled otherwise. You can also enable it with the `--cache` flag or disable it with the `--no-cache` flag.

### Examples

Assume you have `GITHUB_TOKEN` set to `cf4b78a2b8356059f340a7df735d0f63` for the `development` environment in the EnvKey App. You generate a local development `ENVKEY`.

In your project's `.env` file (ignored from source control):

```bash
# .env
ENVKEY=GsL8zC74DWchdpvssa9z-nk7humd7hJmAqNoA
```

Run envkey-source:

```bash
$ eval $(envkey-source)
```

Now `GITHUB_TOKEN` is available in the shell:

```bash
$ echo $GITHUB_TOKEN
cf4b78a2b8356059f340a7df735d0f63
```

Or in any process you launch from this shell:

```bash
$ python
```

```python
>>> import os
>>> os.environ["GITHUB_TOKEN"]
'cf4b78a2b8356059f340a7df735d0f63'
```

You can do exactly the same on a **server**, except instead of putting your `ENVKEY` in a `.env` file, you'll set it as an environment variable (in whatever way you set environment variables for your host/server management platform).

So you set an environment variable on your server:

```bash
ENVKEY=HSyahYDL2jBpyMnkV6gF-2rBFUNAHcQSJTiLA
```

Then you run envkey-source as part of your server start and restart commands, whatever those may be.

```bash
$ eval $(envkey-source) && server-start
```

```bash
$ eval $(envkey-source) && server-restart
```

If you're using envkey-source on a **CI server**, the process is much the same. Set the `ENVKEY` environment variable in your CI interface, then run `eval $(envkey-source)` before running tests.

### Docker

Here's a simple example using Python:

```docker
FROM python:3

# install envkey-source
RUN curl -s https://raw.githubusercontent.com/envkey/envkey-source/master/install.sh | bash

RUN mkdir /code
WORKDIR /code
ADD . /code/

# set EnvKey environment variables before running the process
CMD eval $(envkey-source) && python3 example.py
```

To supply the `ENVKEY` in development with docker-compose, you can add it to a `.env` file, then use the `env_file` key in `docker-compose.yml`.

```yml
services:
  example:
    build: .
    env_file: .env
```

On a server, you just need to pass the ENVKEY environment variable through to your docker container. Where to set this depends on your host, but it shouldn't be difficult.

And now you can access EnvKey variables the same way you'd read normal environment variables.

```python
# example.py

import os

print(os.environ["GITHUB_TOKEN"])
```

### envkey-source within scripts

Note that if you run envkey-source inside a script, your environment variables will only be visible to commands run within that script unless you run the script with `source`, in which case they will be set in the current shell.

### direnv

envkey-source works well with [direnv](https://direnv.net). Just add the following to your `.envrc` file:

```bash
export ENVKEY=HSyahYDL2jBpyMnkV6gF-2rBFUNAHcQSJTiLA

if has envkey-source; then
  eval $(envkey-source --cache)
fi
```

and rerun `direnv allow`.

## x509 error / ca-certificates

On a stripped down OS like Alpine Linux, you may get an `x509: certificate signed by unknown authority` error when `envkey-source` attempts to load your config. [envkey-fetch](https://github.com/envkey/envkey-fetch) (which `envkey-source` wraps) tries to handle this by including its own set of trusted CAs via [gocertifi](https://github.com/certifi/gocertifi), but if you're getting this error anyway, you can fix it by ensuring that the `ca-certificates` dependency is installed. On Alpine you'll want to run:
```
apk add --no-cache ca-certificates
```

## Other EnvKey Libraries

[envkey-fetch](https://github.com/envkey/envkey-fetch) - lower level command line tool that simply accepts an `ENVKEY` and spits out decrypted config as json. Handles core fetching, decryption, verification, web of trust, redundancy, and caching logic. Does most of the work behind the scenes for this library.

[envkey-ruby](https://github.com/envkey/envkey-fetch) - EnvKey Client Library for Ruby and Rails.

[envkey-node](https://github.com/envkey/envkey-node) - EnvKey Client Library for Node.js.

[envkeygo](https://github.com/envkey/envkeygo) - EnvKey Client Library for Go.

## Further Reading

For more on EnvKey in general:

Read the [docs](https://docs.envkey.com).

Read the [integration quickstart](https://docs.envkey.com/integration-quickstart.html).

Read the [security and cryptography overview](https://security.envkey.com).

## Need help? Have questions, feedback, or ideas?

Post an [issue](https://github.com/envkey/envkey-source/issues) or email us: [support@envkey.com](mailto:support@envkey.com).
