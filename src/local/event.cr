module Local
  struct FilePath
    property path : String

    def initialize(@path)
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
  end

  alias TemporalEvent = Start | Stop | Tick | Scanned

  alias Event = FileEvent | TemporalEvent
end
