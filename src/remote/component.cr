require "./cozy"
require "../interfaces/component"

module Remote
  class Component
    include ::ReadyComponent

    def initialize
      @stop = Channel(Nil).new
      @ready = ->{}
      @stopped = ->{}
    end

    def on_ready(&blk)
      @ready = blk
    end

    def on_stopped(&blk)
      @stopped = blk
    end

    def start
      @ready.call
      # TODO: write remote logic
      @stop.receive
    ensure
      @stopped.call
    end

    def stop
      @stop.send nil
    end
  end
end
