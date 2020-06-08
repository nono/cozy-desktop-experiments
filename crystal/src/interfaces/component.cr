module Component
  abstract def start
  abstract def stop
  abstract def on_stopped(&blk)
end

module ReadyComponent
  include Component

  abstract def on_ready(&blk)
end
