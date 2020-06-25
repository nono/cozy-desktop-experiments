# In memory experiment for Cozy desktop

This is an experiment for [Cozy desktop](https://github.com/cozy-labs/cozy-desktop)
where the data is kept in memory (no sqlite/pouchdb). It makes it easier to
test algorithms to synchronize files and directories between the local file
system and the remote Cozy instance. In particular, it is a good way to play
with property based testing and start building a simulator.

## Tooling

I wanted static typing and it looks like OCaml could be a good fit here. I'm
using:

- [Reason](https://reasonml.org/) for the syntax
- [BuckleScript](https://bucklescript.github.io/) as the build tool
- [Belt](https://reasonml.org/apis/javascript/latest/belt) for the stdlib
- [Yarn v1](https://classic.yarnpkg.com/lang/en/) for managing the dependencies
- [Jest](https://jestjs.io/) to run tests
- [Fast-check](https://github.com/dubzzz/fast-check) for property based testing.

You can look at the `scripts` section of the `package.json` to see the useful
commands for developers (the `:w` suffix is used for watch mode, when the task
is runned at each file change).

## OCaml/Reason notes

- 1 file == 1 module (but it is possible to create sub-module inside a file)
- by default, all the `let`s in a module are exported, but an interface file (`.rei`) can restrict that
- it is a good practice to put the signature for exported functions
- the main type of a module is often called `t` (like `Map.String.t`)
- `ref` can be used for a mutable reference, ie a variable which can changes its value
- recursive functions can often be written in a style that allows tail-call optimisation, with an embedded function called by convention `rec` which uses an accumulator called `acc` (also by convention)
- `open` can be used to refer to a content to a module without prefixing it: it is convenient but should be used sparingly as it makes harder to know where some values come from
- don't forget the `;` and the trailing `,`
- file names must be unique per project (no `src/Remote.re` and `src/Fake/Remote.re`)
