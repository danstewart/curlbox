# Curlbox

`curlbox` is a script runner to help make it easier to manage your `curl` scripts and variables.  

See the [tutorial](#tutorial) for a walkthrough of how to use `curlbox`.  

## Install

Grab `curlbox` from the releases and put it in your `PATH`.
<!-- TODO: Add instructions with `go install` -->

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
You should make a new curlbox for each isolated project but separate APIs in the same project can share the same curlbox.

```shell
# Create a new curlbox
curlbox create ~/curlboxes/demo

# Now create a new script inside your curlbox
cat > ~/curlboxes/demo/example.sh << EOF
#!/usr/bin/env bash
curl $URL
EOF

# Make the script executable
chmod +x ~/curlboxes/demo/example.sh

# Now create a vars file
cat > ~/curlboxes/demo/vars.toml << EOF
[default]
URL = "https://example.com"

[some-other-env]
URL = "https://google.com"
EOF

# Now run your new script
# By default the [default] environment will be used
curlbox run ~/curlboxes/demo/example.sh arg1 arg2

# You can specify an environment to use
ENV=some-other-env curlbox run ~/curlboxes/demo/example.sh arg1 arg2

# You can also run in debug mode
DEBUG=1 curlbox run ~/curlboxes/demo/example.sh arg1 arg2
```

Any additional arguments will be passed to your script.  

## Scripts

Scripts can be written in any language but they require the [hashbang](https://en.wikipedia.org/wiki/Shebang_(Unix)) at the top of the file.  
All loaded variables will be accessible in the script as environment variables.

## Variables

Variables are loaded via `vars.toml` and `secrets.toml` files in each directory between the curlbox root and the script directory with each file loaded overriding the previous file.

The `secrets.toml` overrides regular `vars.toml` and are not checked into git by default.

The order the variables are loaded is as follows:
- `vars.toml` at the curlbox root
- `secrets.toml` at the curlbox root
- `vars.toml` in the next directory towards the script
- `secrets.toml` in the next directory towards the script
- etc... until you reach the script directory

```
demo/
├── .curlbox-root
├── pokemon
│   ├── get_by_id.sh
│   ├── secrets.toml     # 4. Inner secrets.toml
│   └── vars.toml        # 3. Inner vars.toml
├── secrets.toml         # 2. Outer most secrets.toml
└── vars.toml            # 1. Outer most vars.toml
```

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

## Recommended Tools

- [curl](https://curl.se/) - A command line tool for making HTTP requests
- [jq](https://stedolan.github.io/jq/) - A command line tool for parsing JSON
- [jless](https://jless.io/) - A command line tool for viewing JSON

## Tips

#### Using formatted JSON via curl

```shell
curl \
    --header "Content-Type: application/json" \
    --request POST \
    --url ${URL} \
    --data-binary @- << EOF
    {
        "name": "$name"
    }
EOF
```

#### Chaining scripts together

`curlbox` will run the script relative to the directory it exists in, this means to chain scripts together you can use the relative path to the script.

```shell
# script1.sh
curl https://example.com
```

```shell
# script2.sh
./script1.sh
```
