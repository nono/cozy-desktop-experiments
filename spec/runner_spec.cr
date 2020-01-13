require "./spec_helper"

class ComponentMock
  include Runner::Component

  property :state

  def initialize
    @mutex = Mutex.new
    @state = :initial
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
    @mutex.synchronize do
      raise "Cannot start from #{@state}" unless [:initial, :stopped].includes? @state
      @state = :starting
    end
    nano = Random.rand 1_000_000
    sleep Time::Span.new(nanoseconds: nano)
    @mutex.synchronize do
      return if @state != :starting
      @state = :ready
    end
    @ready.call
  rescue ex
    p ex
    @state = :error
  end

  def stop
    @mutex.synchronize do
      raise "Cannot start from #{@state}" unless [:starting, :ready].includes? @state
      @state = :stopped
    end
    @stopped.call
  rescue ex
    p ex
    @state = :error
  end
end

describe Runner do
  it "works" do
    local = ComponentMock.new
    remote = ComponentMock.new
    sync = ComponentMock.new
    runner = Runner.new(local, remote, sync)

    5.times do
      spawn do
        nano = Random.rand 10_000_000
        sleep Time::Span.new(nanoseconds: nano)
        runner.stop
        runner.run
      end
    end

    runner.run
    nano = 5_000_000 + Random.rand(10_000_000)
    sleep Time::Span.new(nanoseconds: nano)
    runner.stop
    runner.wait_final_stop

    local.state.should eq :stopped
    remote.state.should eq :stopped
    [:initial, :stopped].should contain sync.state
  end
end
