require "./event"

module Local
  struct Scan
    property path : FilePath

    def initialize(@path)
    end
  end

  struct Close
  end

  struct Checksum
    property path : FilePath

    def initialize(@path)
    end
  end

  alias Effect = Scan | Close | Checksum
end
