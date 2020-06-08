module Simulator
  # MemFS is a mock of a File System where everything is kept in memory.
  #
  # It can be used to simulate a file system without doing syscalls. It uses
  # Time.now, so to avoid the incurred syscalls, it is possible to mock Time
  # with timecop.cr.
  #
  # We don't want to mess the File/Dir classes, and having the possibility to
  # have several MemFS at the same time can be really nice to launch multiple
  # tests in parallel. So, it's better to have a MemFS class that can have
  # several instances, and each instance will be a separated in-memory FS. The
  # drawback is that we have to create an API that is not exactly the same as
  # the File/Dir classes.
  #
  # The mock is limited, not all methods are available (no flock for example),
  # and it can have some differences with a real file system (but we should
  # try to limit the differences as much as possible).
  #
  # The simulated file system is more like Linux than other OS. For example,
  # the filenames are sensible to the case.
  #
  # TODO: add unit tests to check that MemFS behave like a true FS
  class MemFS
    abstract class MemInode
      # TODO: add time and permissions
      property name : String

      # This initialize method is here just to please the compiler
      def initialize(@name : String)
      end
    end

    class MemDir < MemInode
      property children : Array(MemInode)

      def initialize(@name)
        @children = [] of MemInode
      end
    end

    class MemFile < MemInode
      property content : Bytes

      def initialize(@name, @content)
      end
    end

    def initialize
      @entries = [] of MemInode
    end

    def exists?(path)
    end

    def mkdir(path)
    end
  end
end
