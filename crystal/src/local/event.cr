module Local
  # TODO use https://crystal-lang.org/api/0.32.1/Path.html
  record FilePath, path : String do
    # TODO: find a better name for to_s
    def to_s
      @path
    end
  end

  record FileEvent, path : FilePath, type : File::Type
  record Start
  record Stop
  record Tick

  record Scanned, path : FilePath
  record Checksummed, path : FilePath, sum : String

  alias TemporalEvent = Start | Stop | Tick
  alias OperationEvent = Scanned | Checksummed

  # In go, it should probably be something like:
  #
  # type Event interface {
  #   Analyze(store Store) []Effect
  # }
  alias Event = FileEvent | TemporalEvent | OperationEvent
end
