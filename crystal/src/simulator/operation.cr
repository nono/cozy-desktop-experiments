require "json"

module Simulator
  abstract class Operation
    include JSON::Serializable

    # TODO: try to avoid this hard-coded map of types by playing with the hooks
    # https://crystal-lang.org/reference/syntax_and_semantics/macros/hooks.html
    use_json_discriminator "type", {
      start_client: StartClient,
      stop_client:  StopClient,
      sleep:        Sleep,
      create_dir:   CreateDir,
      write_file:   WriteFile,
      move:         Move,
      remove:       Remove,
    }

    property type : String

    def ==(other : Operation)
      self.to_h == other.to_h
    end

    def to_h
      h = {} of Symbol => String | Int32
      {% for ivar in @type.instance_vars.map(&.name) %}
        h[:{{ivar.id}}] = @{{ivar.id}}
      {% end %}
      h
    end
  end

  # TODO: can we use record here?
  # client is the index in the clients list of the scenario
  class StartClient < Operation
    JSON.mapping(
      type: String,
      client: Int32
    )

    @type = "start_client"

    def initialize(@client)
    end
  end

  class StopClient < Operation
    JSON.mapping(
      type: String,
      client: Int32
    )

    @type = "stop_client"

    def initialize(@client)
    end
  end

  class Sleep < Operation
    JSON.mapping(
      type: String,
      ms: Int32
    )

    @type = "sleep"

    def initialize(@ms)
    end
  end

  # ms is the time taken to scan this directory
  # TODO: should we have a field to tell the latency between creating the
  # directory and the inotify event?
  class CreateDir < Operation
    JSON.mapping(
      type: String,
      client: Int32,
      path: String,
      ms: Int32
    )

    @type = "create_dir"

    def initialize(*, @client, @path, @ms)
    end
  end

  # ms is the tuime taken to read the file (e.g. for checksum)
  # TODO: should we have a field to tell the latency between creating the
  # file and the inotify event?
  class WriteFile < Operation
    JSON.mapping(
      type: String,
      client: Int32,
      path: String,
      size: Int32,
      ms: Int32
    )

    @type = "write_file"

    def initialize(*, @client, @path, @size, @ms)
    end
  end

  class Move < Operation
    JSON.mapping(
      type: String,
      client: Int32,
      from: String,
      to: String
    )

    @type = "move"

    def initialize(*, @client, @from, @to)
    end
  end

  class Remove < Operation
    JSON.mapping(
      type: String,
      client: Int32,
      path: String
    )

    @type = "remove"

    def initialize(*, @client, @path)
    end
  end

  # TODO: add more operations like ExcelSave, ToggleNetwork, SelectiveSyncChange
end
