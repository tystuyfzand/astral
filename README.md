Astral
======

A command router for multiple platforms.

Out of the box support currently exists/is focused on [arikawa](https://github.com/diamondburned/arikawa) - a Discord 
bot library.

Importing/Installing
-

```bash
go get github.com/auroradevllc/astral/v3
```

```go
import "github.com/auroradevllc/astral/v3"
```

Signatures
----------

Astral supports signatures, which are a command and arguments defined in a single string.

Example:

```
command <something> <#channel> [optional]
```

This defines a command `command`, with required argument `something`, channel argument `channel`, and optional `optional`.

Middleware
----------

Each route can have middleware assigned to back out/stop execution of a command. This is useful for injecting parameters, checking conditions (Permissions, NSFW), etc.

See the "middleware" folder for examples.

Examples
--------

A basic example showing the usage and middleware is available under `examples/basic`