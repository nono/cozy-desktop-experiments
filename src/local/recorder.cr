require "./analyze"
require "./store"
require "./watcher"

module Local
  # TODO: write documentation for `Recorder`
  class Recorder
    def self.start(dir)
      new(dir).tap &.start
    end

    def initialize(@dir : String)
      @channel = Channel(Event).new(capacity: 1000)
    end

    def start
      store = Store.new
      watcher = Watcher.start(@dir, @channel)

      loop do
        event = @channel.receive
        pp! event
        effects = Local.analyze(store, event)
        pp! effects
        watcher.close if event == TemporalEvent::Stop
      end
    rescue Channel::ClosedError
      # We are done
    end

    def stop
      @channel.send TemporalEvent::Stop
    end
  end
end
