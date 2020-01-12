require "./spec_helper"

class SideMock
  include Runner::Side

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
    local = SideMock.new
    remote = SideMock.new
    sync = SideMock.new
    runner = Runner.new(local, remote, sync)

    spawn runner.run
    5.times do
      spawn do
        nano = Random.rand 100_000_000
        sleep Time::Span.new(nanoseconds: nano)
        runner.stop
        runner.run
      end
    end

    nano = 50_000_000 + Random.rand(100_000_000)
    sleep Time::Span.new(nanoseconds: nano)
    runner.stop
    runner.wait_final_stop

    local.state.should eq :stopped
    remote.state.should eq :stopped
    sync.state.should eq :stopped
  end
end
