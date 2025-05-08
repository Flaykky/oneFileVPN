# one-file VPN

A minimalistic Go-based local SOCKS5 proxy that forwards all traffic through an upstream HTTP proxy. Implemented in a single file for simplicity.

---

## Table of Contents

* [Overview](#overview)
* [Features](#features)
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
* [Command-Line Flags](#command-line-flags)
* [How It Works](#how-it-works)
* [Examples](#examples)
* [Limitations](#limitations)
* [Contributing](#contributing)
* [License](#license)

---

## Overview

This small Go program listens on a local TCP port (binding to `127.0.0.1` on a random free port) and acts as a SOCKS5-compatible proxy. All incoming connections are forwarded through a specified upstream HTTP proxy, optionally using basic authentication.

---

## Features

* Single-file Go implementation
* Zero external dependencies
* Local SOCKS5 proxy on a random port
* HTTP proxy forwarding
* Optional HTTP Basic authentication

---

## Requirements

* Go 1.18 or newer
* Network access to your upstream HTTP proxy

---

## Installation

Clone this repository (or copy the single file) and compile:

```bash
git clone https://example.com/one-file-vpn.git
cd one-file-vpn
go build -o one-file-vpn main.go
```

Alternatively, if you have the `main.go` file only:

```bash
go build -o one-file-vpn main.go
```

---

## Usage

```bash
./one-file-vpn -proxy <host:port> [-auth <username:password>]
```

On startup, the program prints the local listening address:

```
Local SOCKS5 proxy started on 127.0.0.1:53427
```

Configure your applications or system to use a SOCKS5 proxy at that address.

---

## Command-Line Flags

| Flag     | Description                                              | Required | Default |
| -------- | -------------------------------------------------------- | -------- | ------- |
| `-proxy` | Address of the upstream HTTP proxy in `host:port` format | Yes      | —       |
| `-auth`  | Credentials for HTTP Basic auth in `user:password`       | No       | —       |

---

## How It Works

1. **Flag Parsing**

   * The `-proxy` flag is mandatory. If omitted or invalid, the program exits with an error.
   * The optional `-auth` flag must be in `username:password` format; otherwise, the program exits.

2. **Proxy URL Construction**

   * Prepends `http://` to the supplied `host:port`.
   * If credentials are provided, they are embedded in the URL (`http://user:pass@host:port`).

3. **Local Listener**

   * Binds to `127.0.0.1` on a random available port (`:0`).
   * Prints the chosen port for the user to configure their applications.

4. **Connection Handling**

   * For each incoming client connection:

     * Establish a new TCP connection to the upstream HTTP proxy.
     * Relay bytes bidirectionally between client and proxy using `io.Copy`.

---

## Examples

### Without Authentication

```bash
./one-file-vpn -proxy 192.0.2.10:8080
# Output: Local SOCKS5 proxy started on 127.0.0.1:54321
```

Configure your browser’s SOCKS5 proxy to `127.0.0.1:54321`.

### With Authentication

```bash
./one-file-vpn -proxy proxy.example.com:3128 -auth alice:secret
# Output: Local SOCKS5 proxy started on 127.0.0.1:53214
```

---

## Limitations

* **No SOCKS5 handshake**: This implementation does *not* perform the official SOCKS5 handshake; it simply accepts any TCP connection and forwards it. Some clients may require a proper handshake.
* **No DNS proxying**: DNS lookups from the client are not proxied; clients must resolve themselves.
* **Single upstream proxy**: No support for proxy chaining or fallback.
* **No encryption**: Traffic between this proxy and the upstream is plain HTTP CONNECT.

---

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to open an issue or submit a pull request.

---

## License

This project is licensed under the [MIT License](LICENSE).
