# [RFC] A new start for cozy-desktop

## Introduction

[Cozy Drive for Desktop](https://github.com/cozy-labs/cozy-desktop) allows you
to synchronize the files stored in [your Cozy](https://cozy.io) with your
laptop and/or desktop computer. It replicates your files on your hard drive and
applies the changes you make to them on other synchronized devices and on [your
online Cozy](https://github.com/cozy/cozy-stack).

The current version of code kinds of work, but it has accumulated lot of
technical debt, and it will be very hard to go to the next level of reliability
and to feature like selective synchronization (aka "I want to synchronize this
directory but not this one"). This document explain how I will do it if I could
start again. I hope it will be useful someday, and in the mean time, comments
about it are welcomed.


## Limits of the current model

The design of the actual version of cozy-desktop is explained on
https://cozy-labs.github.io/cozy-desktop/doc/developer/design.html.
It was written for Cozy v2, several years ago, with some asumptions that are no
longer relevant. In particular, I made the tradeoff to accept that a move on
the local file system can be sometimes detected and handled as if the file was
deleted and recreated. This is not acceptable for Cozy v3, where a file on the
cozy instance can have references.

I also have more experience on the subject, like knowing the tricks and traps
of inotify, fsevents and ReadDirectoryChangesW. And even if the core
technologies for that haven't changed, there are still some improvements. The
most notable one is that `nodejs` numbers didn't have enough precision for
[fileId on windows](https://github.com/nodejs/node/issues/12115) until very
recently.

But let's talk about the real flaws of the actual design. First and foremost,
it was always obvious to me that having two databases, one for local files and
one for remote files couldn't work and I choose to take a single database
approach with eventual consistency, inspired by CouchDB and PouchDB. I now realize
that having more databases is better: it helps to split the problem in several
sub-problems that are easier to manage. Notably, it's a lot harder to use
inotify, fsevents and ReadDirectoryChangesW that I imagined, and having a
database with local files help to put apart the issues of knowing what happen
on the local file system.

There are also some things that could have been managed better. In particular,
I think it was a mistake to try to resolve conflicts without writing before in
the database. Conflicts are complicated to serialize in the database, but not
doing so is worse: they don't follow the logical flow, duplicate some code, and
introduce a lot of subtle bugs.


## New version

![Workflow](workflow.png)

### Remote

This part is easy: we use the changes feed from CouchDB to put the files
documents in a pouchdb/sqlite/whatever database. In this database, we keep the
relevant data about the directories and files: id and rev, type, name, dir_id,
updated_at, trashed, md5sum, size, executable and tags. But more important, we
don't keep the fullpath (except maybe for debugging purpose).

To know when to pull the documents from the changes feed, we should use the
realtime endpoint from the stack (websocket) to know when there is activity.
With a debouce of 2 seconds, we can be really reactive to the changes made on
the cozy (when online).

This new version is a lot easier than the current one (no need to analyze the
changes feed to regroup the files and directories moved at the same time), even
if we still have to be careful about the transitions between online and offline.

### Local

The local watcher is probably the part that I've spent the most time, and I
still haven't figured how do it in a reliable way. Well, there is the fuse way
where we intercept the syscalls for the files and directories inside the
synchronized path, but it doesn't look like we can do something user friendly
on the 3 platforms (windows, macOS and GNU/Linux, but mainly windows). So let's
forget that. What we want is to have in a database the list of files and
directories, with some fields for each of them. Most of these fields can be
filled with stat(2), with the exception of checksum (a md5sum currently, but it
may be a sha256sum one day).

#### The identifier

The first issue is the identifier to use for this database: the primary key
that will really identify a file or folder, and be used as a foreign key in
other databases. The first obvious candidate is the path. It was my first
choice, in a variant where the path was upper-cased for case-insensitive file
systems. It means that when a folder moves, we have a lot of documents to
update. It can introduce some bugs, for example, if a file in a moved directory
is also updated in the same short period of time. But the more important con
for me is that the same identifier (path) can be used for 2 distinct file or
directory: not at the same time, but one after the other, and if we have a
foreign key from another database, it's really tricky if we can differentiate
if it is the same file or not.

Another candidate is the fileID (windows) or inode number (macOS, GNU/Linux).
It doesn't have the first con (when we move a folder, the identifiers for all
the documents remain the same, we just need to update one document), but it has
the second (the inode numbers can be reused after being freed).

Using an identifier composed from the fileID/inode and a generation has crossed
my mind, but I think that using auto-incremented IDs can work and may be
simpler. We should still have the fileID / inode number in a field (with a
unique index it it's available).

#### Startup

### Backlog

By backlog, I mean the list of files and folders that will need to be
synchronized (or at least to check in case of doubt). In the current version of
cozy-desktop, the backlog is implicit: we keep a sequence number and take the
first document from the pouchdb changes feed after this sequence number. It works
fine, but it has some limitations. We can't query the changes feed or reorder
documents inside it (except by writing to a document that automatically moves
it to the end).

In the new version, I think we should move the backlog to a dedicated database
(or table/keyspace/whatever). We can query the database, and can choose to
synchronize new files before deleting old ones. As an optimization, we can add
QoS by transfering first the small files (< 100 ko), then the medium ones, and
keeps the large ones (> 10Mo) for inactive periods. And, if we want to go
further, we can query a batch of small files and download them in [a single
request as a zip](https://docs.cozy.io/en/cozy-stack/files/#post-filesarchive).

It's also possible to say that, when an error happened, to not retry to
synchronize a file/directory before some amount of time (with an exponential
backoff rule for example).

I don't have a precise list of things to put in the backlog, but we can start
with these properties for each job:
- `side`: `local` or `remote` to say what side has added the job
- `id`: the inode/fileid for `local`, or the uuid if `remote`
- `kind`: `file` or `directory`
- `size`: for files only
- `errors`: the number of errors for synchronzing this file (most often 0)
- `not_before`: we should not try to synchronize the file before this date
  (after an error).

### History

### Sync


## Technologies

This document is mainly focused on a design, how to manage the flow of events
and converge to a stable state. Technologies are not that important for that.
But choosing the right technologies can help to build that faster and with less
bugs. So let's talk a bit about that:

- The language: JavaScript is a terrible language, and even with flow, it's way
  too easy to miss some typos that should be easy to catch bugs. There are hard
  things in cozy-desktop, and it will come with some long sessions of
  debugging, but at least we should avoid that for the bugs that a tool can
  catch (compiler, linter, vet, etc.). We should choose a language that will
  help us to avoid bugs, and JavaScript is not this language (maybe TypeScript
  if we want to keep node and electron, but I haven't written anything
  significant with it, so don't take this as an endorsement).

- The database: PouchDB was not a bad choice, but on the long run, I will
  prefer sqlite. For example, having transactions when writing several documents
  at once after a directory was moved looks more safe. An embedded key/value
  store could also do the job. But I don't feel like we need great performances
  from the database, so the main criterion will be reliability and the ease of
  use (to avoid bugs in misuse). And there, it will be hard to beat sqlite.

- The tests: they are crucial, as the experience have teached us. Mocha is
  mostly fine and have some nice features, but it makes really hard to
  parallelize tests, and I think it's a MUST for testing cozy-desktop. It
  probably means that even if cozy-desktop stays in JS (or something that
  transpiles to it), we should move to another test runner. And I will split
  the tests in 3 categories:

  * unit tests: no cozy-stack, no inotify/FsEvents/RDCW, these tests should be
    run very often by the developers, and it shouldn't take more than a couple
    of seconds to run the whole suite.
  * integration tests: here, we test the boundaries of our code, how the client
    interacts with cozy-stack, with the file system, or with
    inotify/FsEvents/RDCW.
  * global tests: we test the behavior of the client against complex scenarios.
    Fuzzing, property based testing, or just some hand-written scenarios, it
    doesn't matter as long as we are testing a lot of configurations. The goal
    is to have some tools for finding new bugs before the users. But when some
    inputs make the cozy-desktop bugs, it's nice to have a way to keep this
    scenario in the git repository to avoid future regressions.

And that's all!
