require "faker"
require "./scenario"

module Simulator
  class Generate
    property scenario : Scenario

    # TODO: we should add an option for the set of operations and their
    # frequencies
    def initialize(nb_clients, *, @nb_init_ops = 16, @nb_run_ops = 64)
      @scenario = Scenario.new(nb_clients)
    end

    def run
      @scenario.clients.each_with_index do |_, index|
        run_init
        start_client index
      end
      run_ops
      @scenario
    end

    # TODO: add outside path operations
    def run_init
      Random.rand(@nb_init_ops).times do
        freq [{1, ->create_dir}, {1, ->create_file}]
      end
    end

    def run_ops
      # TODO: implement this method
    end

    def freq(choices)
      sum = 0
      choices.each do |choice|
        sum += choice[0]
      end
      nb = Random.rand sum
      action = ->{}
      choices.each do |choice|
        action = choice[1]
        nb -= choice[0]
        break if nb < 0
      end
      action.call
    end

    def start_client(index)
      sleep
      add StartClient.new(client: index)
    end

    def stop_client(index)
      sleep
      add StopClient.new(client: index)
    end

    def sleep
      # TODO: generate more different values
      ms = 2 ** (1 + Random.rand(15))
      add Sleep.new(ms: ms)
    end

    def create_dir
      client = random_client
      path = new_path
      ms = 4 + Random.rand 16
      add CreateDir.new(client: client, path: path, ms: ms)
    end

    def create_file
      client = random_client
      path = new_path
      size = file_size
      ms = 10 + Random.rand(1 + size//1000)
      add CreateFile.new(client: client, path: path, size: size, ms: ms)
    end

    private def random_client
      return 0 if @scenario.clients.size == 1
      Random.rand @scenario.clients.size
    end

    private def new_path
      # TODO: improve it
      Faker::Name.first_name
    end

    private def file_size
      Random.rand(32) + 2 ** (1 + Random.rand(16))
    end

    private def add(op : Operation)
      @scenario.ops << op
    end
  end
end
