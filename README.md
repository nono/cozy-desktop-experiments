# Cozy-fuse

Cozy-fuse is an experiment to use fuse instead of
inotify/FsEvents/ReadDirectoryChangesW for the cozy desktop client. This is a
proof of concept, with a lot of limitations. Don't use for anything else that
some tests, or have backups.

There are other changes between this fuse client and the classical desktop
client:

- Go instead of JS for the progamming langage
- Sqlite3 instead of PouchDB, with several tables
- Some ideas from
  https://github.com/nono/cozy-desktop-experiments/blob/master/RFC/README.md

## What would be needed for a full client

There are a lot of things that are out of the scope for this proof of concept.
This work would be needed if we want to release a new desktop client for Cozy
users:

- Rewrite the code to make it more robust
- Support of Windows and macOS
- UI
- logs and a way to contact the support team
- documentation and tests
- packaging
- use the trash of the local computer
- https://github.com/cozy-labs/cozy-desktop/blob/master/core/config/.cozyignore
- https://github.com/cozy-labs/cozy-desktop/blob/master/core/remote/warning_poller.js
- sentry
- use the XDG user-directories for config & co
- etc.

## How-to test cozy-fuse in local

```sh
$ git clone github.com/nono/cozy-fuse.git
$ cd cozy-fuse
$ go build
$ mkdir -p tmp/{data,mount}
$ cp config.example.json tmp/config.json
$ cozy-stack serve
$ cozy-stack instances add cozy.tools:8080 --passphrase cozy --apps home,store,drive,settings --email foo@cozy.tools --public-name Foo
$ cozy-stack instances client-oauth --json cozy.tools:8080 http://localhost:1234/ 'Cozy Fuse' github.com/nono/cozy-fuse
$ cozy-stack instances token-oauth cozy.tools:8080 <client-id> io.cozy.files
$ cozy-stack instances refresh-token-oauth cozy.tools:8080 <client-id> io.cozy.files
$ $EDITOR tmp/config.json
$ ./cozy-fuse -config tmp/config.json
$ cd tmp/mount && ls -al
```

## License

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

♡2019 by Bruno Michel. Copying is an act of love. Please copy and share.
