require "./analyzer"
require "./watcher"

module Local
  # TODO: write documentation for `Recorder`
  class Recorder
    def self.start(dir)
      new(dir).tap &.start
    end

    def initialize(@dir : String)
      @done = Channel(Nil).new
    end

    def start
      channel = Channel(Event).new
      watcher = Watcher.start(@dir, channel)
      @done.receive?
      watcher.close
    end

    def stop
      @done.send nil
    end
  end
end
