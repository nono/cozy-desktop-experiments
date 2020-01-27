require "./event"

module Local
  struct Scan
    property path : FilePath

    def initialize(@path)
    end
  end

  struct BeReady
  end

  struct Close
  end

  struct ComputeChecksum
    property path : FilePath

    def initialize(@path)
    end
  end

  alias Effect = Scan | BeReady | Close | ComputeChecksum
end
