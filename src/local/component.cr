require "./analyze"
require "./store"
require "./watcher"

module Local
  class Component
    def initialize(@dir : String)
      @channel = Channel(Event).new(capacity: 1000)
      @event = ->{}
      @ready = ->{}
      @stopped = ->{}
    end

    def on_event(&blk)
      @event = blk
    end

    def on_ready(&blk)
      @ready = blk
    end

    def on_stopped(&blk)
      @stopped = blk
    end

    def start
      store = Store.new
      watcher = Watcher.start(@dir, @channel)

      loop do
        event = @channel.receive
        @event.call event
        effects = Local.analyze(store, event)
        effects.each do |effect|
          @ready.call if effect.is_a? BeReady
          watcher.apply effect
        end
      end
    rescue Channel::ClosedError
      # We are done
      @stopped.call
    end

    def stop
      @channel.send TemporalEvent::Stop
    end
  end
end
