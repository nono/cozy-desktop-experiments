require "faker"
require "./scenario"

module Simulator
  class Generate
    SPECIAL_CHARS = [':', '-', 'é', ' ', '%', ',', '&', '@', 'É', 'Ç']

    property scenario : Scenario

    # TODO: we should add an option for the set of operations and their
    # frequencies
    def initialize(nb_clients, *, @nb_init_ops = 16, @nb_run_ops = 64)
      @scenario = Scenario.new(nb_clients)
      @known_paths = [] of String
      @deleted_paths = [] of String
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

    # TODO: find a way to generate operations for files in conflict (we don't
    # known their name/path when generating the scenario, so, we should
    # probably use something like a number to identify them).
    def run_ops
      Random.rand(@nb_run_ops).times do
        freq [
          {5, ->create_dir},
          {3, ->create_file},
          {1, ->recreate_deleted_dir},
          {1, ->recreate_deleted_file},
          {1, ->update_file},
          {2, ->move_to_new_path},
          {3, ->move_to_deleted_path},
          # TODO: add outside operations
          # {1, ->move_to_outside},
          # {1, ->move_from_outside},
          {5, ->remove},
          {2, ->sleep},
        ]
      end
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

    # Operations

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

    def create_file
      write_file new_path
    end

    def create_dir
      create_directory new_path
    end

    def recreate_deleted_dir
      create_directory deleted_path
    end

    def recreate_deleted_file
      write_file deleted_path
    end

    def update_file
      write_file known_path
    end

    def create_directory(path)
      client = random_client
      ms = 4 + Random.rand 16
      add CreateDir.new(client: client, path: path, ms: ms)
    end

    def write_file(path)
      client = random_client
      size = file_size
      ms = 10 + Random.rand(1 + size//1000)
      add WriteFile.new(client: client, path: path, size: size, ms: ms)
    end

    def move_to_new_path
      from = known_path
      to = new_path
      move from, to
    end

    def move_to_deleted_path
      from = known_path
      to = deleted_path
      move from, to
    end

    def move(from, to)
      @deleted_paths << from
      client = random_client
      add Move.new(client: client, from: from, to: to)
    end

    def remove
      client = random_client
      path = known_path
      add Remove.new(client: client, path: path)
    end

    # Low-level helpers

    private def random_client
      return 0 if @scenario.clients.size == 1
      Random.rand @scenario.clients.size
    end

    private def known_path
      return new_path if @known_paths.empty?
      @known_paths.sample
    end

    private def deleted_path
      return new_path if @deleted_paths.empty?
      @deleted_paths.sample
    end

    private def new_path
      # TODO: improve it
      p = Faker::Name.first_name
      p = SPECIAL_CHARS.sample + p if Random.rand(4).zero?
      p = @known_paths.sample + "/" + p unless @known_paths.empty? || Random.rand(3).zero?
      @known_paths << p
      p
    end

    private def file_size
      Random.rand(32) + 2 ** (1 + Random.rand(16))
    end

    private def add(op : Operation)
      @scenario.ops << op
    end
  end
end
