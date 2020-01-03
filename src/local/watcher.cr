require "inotify"

module Local
  # TODO: Write documentation for `Watcher`
  class Watcher
    def self.start(dir, channel)
      new(dir, channel)
    end

    property dir : String
    property channel : Channel(Symbol)
    property inotify : Inotify::Watcher

    def initialize(@dir, @channel)
      @inotify = spawn_inotify
      spawn_ticker
    end

    private def spawn_ticker
      spawn do
        sleep seconds: 0.1
        @channel.send :tick
      end
    end

    private def spawn_inotify
      Inotify.watch @dir do |event|
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
