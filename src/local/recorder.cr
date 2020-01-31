require "./component"

module Local
  # TODO: write documentation for `Recorder`
  class Recorder
    def initialize(dir : String)
      @component = Component.new(dir)
      @component.on_event { |event| on_event(event) }
      @ticks = 0
    end

    def start
      @component.start
    end

    def stop
      @component.stop
    end

    def on_event(event : Tick)
      @ticks += 1
    end

    def on_event(event : Start)
      print_ticks
      puts "Start"
    end

    def on_event(event : Stop)
      print_ticks
      puts "Stop"
    end

    def on_event(event : FileEvent | OperationEvent)
      print_ticks
      pp! event
    end

    def print_ticks
      puts "+#{@ticks}" if @ticks > 0
      @ticks = 0
    end
  end
end
