> This is a fork from the awesome project from @maaslalani with some modifications.

> The main change here is generate invoice for services and not products, so you can use this to generate invoice for your freelance or contractor work.

# Invoice

Generate invoices from the command line.

## Command Line Interface

```bash
invoice generate --from "Dream, Inc." --to "Imagine, Inc." \
    --service "Software development" --rate 2500 \
    --note "For debugging purposes."
```

View the generated PDF at `invoice.pdf`, you can customize the output location
with `--output`.

```bash
open invoice.pdf
```

### Configuration File

Or, save repeated information with JSON / YAML:

```json
{
    "logo": "/path/to/image.png",
    "from": "Dream, Inc.",
    "to": "Imagine, Inc.",
    "services": ["Software development"],
    "rates": [2500],
}
```

Generate new invoice by importing the configuration file:

```bash
invoice generate --import path/to/data.json \
    --output duck-invoice.pdf
```

## Installation

Install with Go:

```sh
go install github.com/wesleimp/invoice@main
```

Or download a binary from the [releases](https://github.com/wesleimp/invoice/releases).

## License

[MIT](https://github.com/wesleimp/invoice/blob/master/LICENSE)

<sub><sub>z</sub></sub><sub>z</sub>z
