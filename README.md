Astral
======

A command router for [discordgo](https://github.com/bwmarrin/discordgo) with a few twists.

Heavily inspired by [dgrouter](https://github.com/Necroforger/dgrouter), but based off the command system used in [Astra](https://astrabot.net).

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

