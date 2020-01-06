module Local
  struct FilePath
    property path : String

    def initialize(@path)
    end
  end

  class FileEvent
    property path : FilePath

    def initialize(@path)
    end
  end

  enum TemporalEvent
    Start
    Stop
    Tick
  end

  alias Event = FileEvent | TemporalEvent
end
