struct Time
  def after(time : Time::Span)
    channel = Channel(Nil).new(1)
    spawn do
      sleep time
      channel.send nil
    end
    channel
  end
end
