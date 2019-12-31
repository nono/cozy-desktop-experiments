require "inotify"

module Local
  # TODO: Write documentation for `Watcher`
  class Watcher
    def self.start(dir)
      new(dir)
    end

    def initialize(dir)
      @inotify = Inotify.watch dir do |event|
        pp! event
        path = File.join([event.path, event.name].select(String))
        begin
          inode = File.info(path).ino
          pp! inode
        rescue
          pp "no inode number"
        end
      end
    end

    def close
      @inotify.close
    end
  end
end
