# Support the 3 major OS: GNU/Linux, macOS, and Windows

## Naming files

- https://en.wikipedia.org/wiki/Filename
- https://docs.microsoft.com/fr-fr/windows/win32/fileio/naming-a-file
- https://mjtsai.com/blog/2017/06/27/apfs-native-normalization/
- http://www.wsanchez.net/papers/USENIX_2000/
- https://help.dropbox.com/en-us/installs-integrations/sync-uploads/files-not-syncing

## Extended attributes

GNU/Linux and macOS have [extended file attributes](https://en.wikipedia.org/wiki/Extended_file_attributes).
On windows, alternate data streams can be used to emulate them.

- https://www.freedesktop.org/wiki/CommonExtendedAttributes/

### Why not using extended attributes to track files moved/renamed?

1. The hard case is a software that write a new version of a file in a
   temporary place, and then replaces the file but the new version. The
   extended attributes doesn't help us in this case. So, if an event from the
   FS watcher is for file with no extended attributes, we still have to rely
   on another way to check if it is the same file or not than a previous file.
2. And, in the other case, for an event for a file with an extended attribute,
   we have to check if it is a move/rename or a copy. And, again, we have to
   rely on another method to do check that.
3. Reading the extended attributes is a system call, and as such, it takes
   time, and during this time, other operations can take place on the file
   system. Adding latency to the analyze increases the risk of race conditions.

So, using extended attributes has a high cost and a low signal for this use
case. And that's why we don't use them. But, I still think they can be useful
in specific cases:

- on Linux, when the client is stopped, FS operations can lead to inodes being
  reused
- when a desktop client is revoked and connected again.

In those two cases, information has been lost by the client, and it can't use
its usual method to detect files that have been moved/renamed in a reliable way.
Extended file attributes could be useful here.

## Timestamps for the file systems

### Windows

Windows has 3 timestamps for each file or directory:

- CreationTime
- LastAccessTime
- LastWriteTime

See https://docs.microsoft.com/en-us/windows/desktop/api/FileAPI/nf-fileapi-getfiletime
and https://docs.microsoft.com/fr-fr/windows/desktop/api/fileapi/nf-fileapi-setfiletime

**Limitations**

The three timestamps can be read and changed, but nodejs currently does not
support changing the CreationTime.

The resolution for each timestamps can be different:

> For example, on FAT, create time has a resolution of 10 milliseconds, write
> time has a resolution of 2 seconds, and access time has a resolution of 1 day 

### Linux

On Linux, a file can have 4 timestamps:

- atime, the time of the last access (read or write)
- ctime, the time of the last change (it can be the content or the metadata like in a chmod case)
- mtime, the time of the last modification of the content
- birthtime (also called crtime), the time of the creation of the file

See http://man7.org/linux/man-pages/man2/stat.2.html
and http://man7.org/linux/man-pages/man2/statx.2.html

**Limitations**

The resolution for timestamps is not same for all the file systems.

Birthtime is not available on some file systems.

Atime and mtime can be changed easily. Ctime can only be changed to the current
time. Birthtime cannot be changed.

Updating atime at each access costs a lot of performances, so the file systems
are usually mounted with some options to avoid that (like `noatime` or
`relatime`).

No order on the timestamp should be presumed. It's possible for a file to have:
`atime < ctime < mtime < birthtime` (even if it is generally in the reverse order).

### macOS

On macOS, a file has the same 4 timestamps that on linux: atime, ctime, mtime, and birthtime.

See https://www.sciencepubco.com/index.php/ijet/article/view/13870

**Limitations**

The resolution for timestamps is the second for HFS+ and the nanosecond for APFS.

The atime and mtime can be easily changed. For the birthtime, there is a tool
called [`SetFile`](https://www.unix.com/man-page/osx/1/SetFile/) that can do
that, and touch will also modify the birthtime when trying to set a new mtime
that is before the birthtime with the `-t` option. Ctime can only be set to the
current time.

