require "inotify"
require "./event"

module Local
  # TODO: Write documentation for `Watcher`
  class Watcher
    def self.start(dir, channel)
      new(dir, channel)
    end

    property dir : String
    property channel : Channel(Event)
    property inotify : Inotify::Watcher

    def initialize(@dir, @channel)
      @inotify = inotify
      spawn ticker
      @channel.send TemporalEvent::Start
    end

    def ticker
      spawn do
        loop do
          sleep seconds: 0.1
          @channel.send TemporalEvent::Tick
        end
      end
    end

    def inotify
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
