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

  class FileEvent
    property path : FilePath
    property type : File::Type

    def initialize(@path, @type)
    end
  end

  struct Start
  end

  struct Stop
  end

  struct Tick
  end

  struct Scanned
    property path : FilePath

    def initialize(@path)
    end
  end

  struct Checksummed
    property path : FilePath
    property sum : String

    def initialize(@path, @sum)
    end
  end

  alias TemporalEvent = Start | Stop | Tick

  alias OperationEvent = Scanned | Checksummed

  # In go, it should probably be something like:
  #
  # type Event interface {
  #   Analyze(store Store) []Effect
  # }
  alias Event = FileEvent | TemporalEvent | OperationEvent
end
