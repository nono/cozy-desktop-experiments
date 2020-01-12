# The runner ensures the coordination between 3 components: the local and
# remote sides, and sync. It is expected that the runner is created and then is
# called by this methods:
# - first, a run to start working when the desktop client starts
# - then, fibers can make a restart by calling stop and then run
# - and finally, the stop method can be called before the client exits.
#
# Constraints:
# - the local and remote sides must be started in parallel
# - the sync must be started when both the local and remote sides are ready
# - when stop is called on the runner, the stop method must be called on the
#   started components
# - the components can take some time to stop
# - if a few restarts happen in a short time, we need to ensure that all
#   components are stopped before trying to start them again.
# - _bonus_ if the runner know that a run should be stopped, it can avoid
#   starting the components and stopping them just after.
class Runner
  enum Message
    Ready
    Stop
    Stopped
  end

  module Side
    abstract def start
    abstract def stop
    abstract def on_ready(&blk)
    abstract def on_stopped(&blk)
  end

  def initialize(@local : Side, @remote : Side, @sync : Side)
    @chan = Channel(Message).new(3)
    @can_start = Mutex.new
    @can_stop = Channel(Symbol).new(1)
    @count = 0

    @local.on_ready { @chan.send Message::Ready }
    @remote.on_ready { @chan.send Message::Ready }

    @local.on_stopped { @chan.send Message::Stopped }
    @remote.on_stopped { @chan.send Message::Stopped }
    @sync.on_stopped { @chan.send Message::Stopped }
  end

  # run is a blocking call to starts the 3 components with the good time, and
  # to stop them when asked for it.
  def run
    @can_start.lock
    @can_stop.send :token

    if !cancelled?
      spawn @local.start
      spawn @remote.start
      @count = 2

      if initialized
        spawn @sync.start
        @count += 1
        wait_for_stop
        spawn @sync.stop
      end

      spawn @remote.stop
      spawn @local.stop

      wait_all_stopped
    end
  ensure
    @can_start.unlock
  end

  private def cancelled?
    timeout = after(Time::Span.new(nanoseconds: 100_000))
    received = Channel.receive_first(@chan, timeout)
    return received == Message::Stop
  end

  private def initialized
    2.times do
      received = @chan.receive
      return false if received == Message::Stop
    end
    return true
  end

  private def wait_for_stop
    loop do
      msg = @chan.receive
      return if msg == Message::Stop
    end
  end

  private def wait_all_stopped
    while @count > 0
      msg = @chan.receive
      @count -= 1 if msg == Message::Stopped
    end
  end

  # stop can be called to stop the current run. It can blocks if several
  # restarts happen in a short time. So, it is probably better to make the
  # restarts in their own fiber:
  #
  #    spawn do
  #      runner.stop
  #      runner.run
  #    end
  def stop
    @can_stop.receive
    @chan.send Message::Stop
  end

  # wait_final_stop can be called to ensure that the runner is stopped and
  # cannot restart. It is useful when the client wants to exit.
  def wait_final_stop
    @can_start.lock
  end

  # TODO: find another place to put this method
  private def after(time : Time::Span)
    channel = Channel(Nil).new(1)
    spawn do
      sleep time
      channel.send nil
    end
    channel
  end
end
