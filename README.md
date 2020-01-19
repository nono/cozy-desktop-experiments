# cozy-desktop.cr

Try Crystal programming language by writing a desktop client for Cozy Cloud
that synchronizes files.

This is a proof of concept, with a lot of limitations. Don't use for anything
else that some tests, or have backups.

## What would be needed for a full client

There are a lot of things that are out of the scope for this proof of concept.
This work would be needed if we want to release a new desktop client for Cozy
users:

- Rewrite the code to make it more robust
- Support of Windows and macOS
- UI
- logs and a way to contact the support team
- fix the TODOs in code
- more tests (a lot, manual and automatic)
- packaging & auto-update
- auto-start
- use the trash of the local computer
- https://github.com/cozy-labs/cozy-desktop/blob/master/core/config/.cozyignore
- https://github.com/cozy-labs/cozy-desktop/blob/master/core/remote/warning_poller.js
- sentry
- use the XDG user-directories for config & co
- use a specific User-Agent
- update the version of the registered OAuth2 client
- call the `/settings/synchronized` endpoint
- managing correctly the several dates for each file
- handling `.cozy-note` files
- updating the dependencies with renovate or dependabot
- write more documentation and publish it (`crystal docs`) on GitHub pages
- etc.

## Links

* https://crystal-lang.org/2019/04/30/watch-run-change-build-repeat.html
* https://crystal-lang.org/2018/09/04/using-circleci-2.0-for-your-crystal-projects.html
* https://github.com/crystal-ameba/ameba
* https://github.com/ysbaddaden/http2
* https://github.com/petoem/inotify.cr
* https://github.com/waterlink/quick.cr
* https://github.com/crystal-community/hardware
* https://github.com/Sija/retriable.cr
* https://github.com/chris-huxtable/atomic_write.cr
* https://github.com/spalger/crystal-mime
* https://github.com/crystal-lang/crystal-sqlite3
* https://github.com/TPei/progress_bar.cr
* https://github.com/blacksmoke16/crylog
* https://github.com/ysbaddaden/earl
* https://quicktype.io/

## Installation

You will need to [install Crystal](https://crystal-lang.org/install/) and run:

```sh
$ git clone git@github.com:nono/cozy-desktop.cr.git
$ cd cozy-desktop.cr
$ shards install
$ shards build
$ ./bin/cozy-desktop-ng --help
```

## Usage

You can use these commands to play with a local cozy instance:

```sh
$ cd cozy-desktop.cr
$ mkdir -p tmp
$ cozy-stack serve
$ cozy-stack instances add cozy.tools:8080 --passphrase cozy --apps home,store,drive,settings --email foo@cozy.tools --public-name Foo
$ cozy-stack instances client-oauth --json cozy.tools:8080 http://localhost:1234/ 'Cozy Fuse' github.com/nono/cozy-fuse
$ cozy-stack instances token-oauth cozy.tools:8080 <client-id> io.cozy.files
$ cozy-stack instances refresh-token-oauth cozy.tools:8080 <client-id> io.cozy.files
$ ./bin/cozy-desktop-ng configure --dir ./tmp/Cozy --cozy cozy.tools:8080 --token <token>
$ ./bin/cozy-desktop-ng sync
$ ls -alr ./tmp/Cozy
```

## Development

Some useful commands:

- `crystal spec` will run the unit tests
- `crystal tool format --check` will check that the `.cr` files are correctly
  formatted
- `crystal docs` will build the documentation inside the `docs/` directory
- `shards update` can be used to update the dependencies

## Contributing

1. Fork it (<https://github.com/nono/cozy-desktop.cr/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## Contributors

- [Bruno Michel](https://github.com/nono) - creator and maintainer

## License

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

♡2020 by Bruno Michel. Copying is an act of love. Please copy and share.
