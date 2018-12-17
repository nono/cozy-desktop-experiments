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

