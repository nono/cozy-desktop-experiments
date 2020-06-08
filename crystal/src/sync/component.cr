require "../interfaces/component"

module Sync
  class Component
    include ::Component

    def initialize
      @stop = Channel(Nil).new
      @stopped = ->{}
    end

    def on_ready(&blk)
    end

    def on_stopped(&blk)
      @stopped = blk
    end

    def start
      # TODO: write sync logic
      @stop.receive
    ensure
      @stopped.call
    end

    def stop
      @stop.send nil
    end
  end
end
