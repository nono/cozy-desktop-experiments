require "./event"

module Local
  struct Scan
    def initialize(@path : FilePath)
    end

    def path
      @path.path
    end
  end

  struct BeReady
  end

  struct Close
  end

  struct Checksum
    def initialize(@path : FilePath)
    end

    def path
      @path.path
    end
  end

  alias Effect = Scan | BeReady | Close | Checksum
end
