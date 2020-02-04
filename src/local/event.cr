module Local
  # TODO use https://crystal-lang.org/api/0.32.1/Path.html
  struct FilePath
    property path : String

    def initialize(@path)
    end

    # TODO: find a better name for to_s
    def to_s
      @path
    end
  end

  alias FileEvent = NamedTuple(path: FilePath, type: File::Type)

  struct Start
  end

  struct Stop
  end

  struct Tick
  end

  alias Scanned = NamedTuple(path: FilePath)
  alias Checksummed = NamedTuple(path: FilePath, sum: String)

  alias TemporalEvent = Start | Stop | Tick
  alias OperationEvent = Scanned | Checksummed

  # In go, it should probably be something like:
  #
  # type Event interface {
  #   Analyze(store Store) []Effect
  # }
  alias Event = FileEvent | TemporalEvent | OperationEvent
end
