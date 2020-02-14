require "./event"

module Local
  record Scan, path : FilePath
  record BeReady
  record Close
  record ComputeChecksum, path : FilePath

  alias Effect = Scan | BeReady | Close | ComputeChecksum
end
