module Simulator
  struct StartClient
    property nb : Int32

    def initialize(@nb)
    end
  end

  struct StopClient
    property nb : Int32

    def initialize(@nb)
    end
  end

  struct Sleep
    property ms : Int32

    def initialize(@ms)
    end
  end

  struct CreateDir
    property client : Int32
    property path : String
    property ms : Int32 # How much time does it take to scan this dir

    def initialize(*, @client, @path, @ms)
    end
  end

  struct CreateFile
    property client : Int32
    property path : String
    property size : Int32
    property ms : Int32 # How much time does it take to read this file

    def initialize(*, @client, @path, @size, @ms)
    end
  end

  alias Operation = StartClient |
                    StopClient |
                    Sleep |
                    CreateFile |
                    CreateDir
  # TODO: add more operations like ToggleNetwork
end
