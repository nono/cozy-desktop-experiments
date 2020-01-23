require "./component"

module Local
  # TODO: write documentation for `Recorder`
  class Recorder
    def self.start(dir)
      new(dir).tap &.start
    end

    def initialize(dir : String)
      @component = Component.new(dir)
      # @component.on_event { |event| on_event(event) }
    end

    def on_event(event)
      ticks = 0
      if event == TemporalEvent::Tick
        ticks += 1
      else
        puts "+#{ticks}" if ticks > 0
        ticks = 0
        pp! event
      end
    end

    def stop
      @component.stop
    end
  end
end
