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
      ticks = 0

      loop do
        event = @channel.receive
        if event == TemporalEvent::Tick
          ticks += 1
        else
          puts "+#{ticks}" if ticks > 0
          ticks = 0
          pp! event
        end
        effects = Local.analyze(store, event)
        effects.each do |effect|
          pp! effect
          watcher.apply effect
        end
      end
    rescue Channel::ClosedError
      # We are done
    end

    def stop
      @channel.send TemporalEvent::Stop
    end
  end
end
