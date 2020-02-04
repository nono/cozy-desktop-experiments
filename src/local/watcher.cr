require "inotify"
require "digest/md5"
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
      @channel.send Start.new
      spawn ticker
    end

    def ticker
      spawn do
        loop do
          sleep seconds: 0.1
          @channel.send Tick.new
        end
      end
    end

    def apply(effect : Scan)
      dir = File.join @dir, effect.path.to_s
      pp! dir
      # TODO: check what happens if the dir is deleted just right when we try
      # to watch it
      notifier = Inotify.watch dir do |event|
        pp! event
        fullpath = File.join([event.path, event.name].select(String))
        prepare_file_event fullpath
      end
      @notifiers << notifier
      Dir.entries(dir).each do |name|
        next if name == "." || name == ".."
        fullpath = File.join dir, name
        prepare_file_event fullpath
      end
      @channel.send Scanned.new(effect.path)
    end

    def apply(effect : Close)
      @notifiers.each &.close
      @notifiers = [] of Inotify::Watcher
    end

    def apply(effect : ComputeChecksum)
      Fiber.new do
        fullpath = File.join @dir, effect[:path].to_s
        sum = Digest::MD5.hexdigest do |ctx|
          File.open fullpath do |f|
            slice = Bytes.new(4096)
            while f.read(slice) > 0
              ctx.update slice
            end
          end
        end
        event = Checksummed.new(path: effect[:path], sum: sum)
        @channel.send event
      end
    end

    def apply(effect : BeReady)
    end

    def prepare_file_event(fullpath : String)
      info = File.info(fullpath)
      pp! info.ino
      path = FilePath.new(fullpath.lchop @dir)
      @channel.send FileEvent.new(path, info.type)
    rescue
      pp "no inode number"
    end
  end
end
