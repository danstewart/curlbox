# Curlbox

A box of curl commands.

## Install

Grab `curlbox` from the releases and put it in your `PATH`.

## Usage Reference

```shell
# Create a new curlbox
curlbox create path/to/store/curlbox

# Run a script
# Any additional arguments will be passed to the script
curlbox run path/to/script [args...]

# Or run in debug mode
DEBUG=1 curlbox run path/to/script [args...]
```

## Tutorial

A curlbox is just a directory containing variables and scripts for making HTTP requests.  
You should make a new curlbox for each isolate project but separate APIs in the same project can share the same curlbox.

```shell
# Create a new curlbox
curlbox create ~/curlboxes/demo

# Now create a new script inside your new curlbox:
echo "curl $URL" > ~/curlboxes/demo/example.sh

# Now create a vars file:
echo -e "[default]\nURL=http://example.com" > ~/curlboxes/demo/vars.toml

# Now run your new script:
curlbox run ~/curlboxes/demo/example.sh

# You can also run in debug mode using:
DEBUG=1 curlbox run ~/curlboxes/demo/example.sh
```

Any additional arguments will be passed to your script.  

## Vars

Variables are loaded via the `vars.toml` and `secrets.toml` files in each directory between the curlbox root and the script directory with each directory closer to the script overriding the higher up directories.  

The `secrets.toml` override regular `vars.toml` and are not checked into git by default.

## Scripts

Scripts can be written in any language but they require the hashbang at the top of the file.  
All vars will be accessible in the script as environment variables.

## Environments

The `vars.toml` and `secrets.toml` files can be used to define different environments, the default environment is `default`.

```shell
[default]
URL="http://localhost:1234"

[dev]
URL="http://dev.example.com"

[prod]
URL="http://example.com"
```

Then run with:
```shell
ENV=dev curlbox run path/to/script
```

