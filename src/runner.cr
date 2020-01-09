class Runner
  enum Message
    Ready
    Stop
    Stopped
  end

  def initialize(@local, @remote, @sync)
    @chan = Channel(Message).new(3)
    @can_start = Channel(Atom).new(1)
    @can_stop = Channel(Atom).new(1)
    @count = 0

    @local.on_ready { @chan.send Ready }
    @remote.on_ready { @chan.send Ready }

    @local.on_stopped { @chan.send Stopped }
    @remote.on_stopped { @chan.send Stopped }
    @sync.on_stopped { @chan.send Stopped }

    @can_start.send :token
  end

  def run
    @can_start.receive
    @can_stop.send :token

    spawn @local.start
    spawn @remote.start
    @count = 2

    if !cancelled? && initialized
      spawn @sync.start
      @count += 1
      wait_for_stop
      spawn @sync.stop
    end

    spawn @remote.stop
    spawn @local.stop

    wait_all_stopped

    @can_start.send :token
  end

  def cancelled?
    timeout = Time.after(Time::Span.new(nanoseconds: 100_000))
    received = Channel.receive_first(@chan, timeout)
    return received == Stop
  end

  def initialized
    2.times do
      received = @chan.receive
      return false if received == Stop
    end
    return true
  end

  def wait_for_stop
    loop do
      msg = @chan.receive
      return if msg.what == Stop
    end
  end

  def wait_all_stopped
    while @count > 0
      msg = @chan.receive
      @count -= 1 if msg.what == Stopped
    end
  end

  def stop
    @can_stop.receive
    @chan.send Stop
  end
end
