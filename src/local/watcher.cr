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
    property notifiers : Array(Inotify::Watcher)

    def initialize(@dir, @channel)
      @notifiers = [] of Inotify::Watcher
      @channel.send TemporalEvent::Start
      spawn ticker
    end

    def ticker
      spawn do
        loop do
          sleep seconds: 0.1
          @channel.send TemporalEvent::Tick
        end
      end
    end

    def apply(effect : Scan)
      dir = File.join @dir, effect.path
      pp! dir
      notifier = Inotify.watch dir do |event|
        pp! event
        path = File.join([event.path, event.name].select(String))
        begin
          info = File.info(path)
          pp! info.ino
          path = FilePath.new(path.lchop @dir)
          @channel.send FileEvent.new(path)
        rescue
          pp "no inode number"
        end
      end
      @notifiers << notifier
    end

    def apply(effect : Close)
      @notifiers.each &.close
      @notifiers = [] of Inotify::Watcher
    end

    def apply(effect : Checksum)
      # TODO
    end
  end
end
