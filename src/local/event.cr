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

  enum TemporalEvent
    Start
    Stop
    Tick
    Scanned
  end

  alias Event = FileEvent | TemporalEvent
end
