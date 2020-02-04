module Simulator
  # client is the index in the clients list of the scenario
  alias StartClient = NamedTuple(client: Int32)
  alias StopClient = NamedTuple(client: Int32)
  alias Sleep = NamedTuple(ms: Int32)

  # ms is the time taken to scan this directory
  alias CreateDir = NamedTuple(client: Int32, path: String, ms: Int32)

  # ms is the tuime taken to read the file (e.g. for checksum)
  alias CreateFile = NamedTuple(client: Int32, path: String, size: Int32, ms: Int32)

  alias Operation = StartClient |
                    StopClient |
                    Sleep |
                    CreateFile |
                    CreateDir
  # TODO: add more operations like ToggleNetwork
end
