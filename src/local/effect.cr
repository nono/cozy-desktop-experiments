require "./event"

module Local
  alias Scan = NamedTuple(path: FilePath)

  struct BeReady
  end

  struct Close
  end

  alias ComputeChecksum = NamedTuple(path: FilePath)

  alias Effect = Scan | BeReady | Close | ComputeChecksum
end
