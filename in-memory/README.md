# In memory experiment for Cozy desktop

This is an experiment for [Cozy desktop](https://github.com/cozy-labs/cozy-desktop)
where the data is kept in memory (no sqlite/pouchdb). It makes it easier to
test algorithms to synchronize files and directories between the local file
system and the remote Cozy instance. In particular, it is a good way to play
with property based testing and start building a simulator.

I wanted static typing and it looks like OCaml could be a good fit here. I'm
using:

- [Reason](https://reasonml.org/) for the syntax
- [BuckleScript](https://bucklescript.github.io/) as the build tool
- [Belt](https://reasonml.org/apis/javascript/latest/belt) for the stdlib
- [Yarn v1](https://classic.yarnpkg.com/lang/en/) for managing the dependencies
- [Jest](https://jestjs.io/) to run tests

You can look at the `scripts` section of the `package.json` to see the useful
commands for developers.
