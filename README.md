# cro-load-test

`cro-load-test` is a distributed load testing tool (and framework) for load
testing [Crypto.com Crossfire](https://chain.crypto.com/crossfire) networks based on [tm-load-test](https://github.com/informalsystems/tm-load-test).

## Requirements
In order to build and use the tools, you will need:

* Go 1.13+
* `make`

## Building
To build the `cro-load-test` binary in the `build` directory:

```bash
make build
```

## Usage
`cro-load-test` can be executed in **standalone** mode

### Standalone Mode
In standalone mode, `cro-load-test` operates in a similar way to `tm-bench`:

```bash
cro-load-test -b 5 -r 1000 \
    --broadcast-tx-method async \
    --endpoints ws://127.0.0.1:26657/websocket \
    --gas 200000 --gas-prices 0.2basetcro
```

To see a description of what all of the parameters mean, simply run:

```bash
cro-load-test --help
```

## Development
To run the linter and the tests:

```bash
make lint
make test
```

