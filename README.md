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
- etc.

## License

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

♡2019 by Bruno Michel. Copying is an act of love. Please copy and share.
