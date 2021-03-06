require "../local/*"
require "../remote/*"
require "../sync/*"
require "./memfs"

module Simulator
  class Play
    record Event, at : Int32, op : Operation
    record Client, fs : MemFS, local : Local::Store

    def initialize(@scenario : Scenario)
      @now = 0
      @events = [] of Event
      @clients = [] of Client
      @scenario.clients.each do
        fs = MemFS.new
        local = Local::Store.new
        @clients << Client.new(fs: fs, local: local)
      end
    end

    def run
    end

    def check
      # TODO: check that the all the clients have the same files and dirs that
      # the cozy instance
      # TODO: can we check other properties? Maybe that files don't disappear
      [] of String
    end
  end
end
