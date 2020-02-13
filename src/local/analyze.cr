require "./effect"
require "./event"
require "./store"

module Local
  extend self

  ROOT = FilePath.new("/")

  def analyze(store : Store, event : Start) : Array(Effect)
    store.scan_counter += 1
    [Scan.new(path: ROOT)] of Effect
  end

  def analyze(store : Store, event : Stop) : Array(Effect)
    [Close.new] of Effect
  end

  def analyze(store : Store, event : Tick) : Array(Effect)
    # Nothing for the moment
    [] of Effect
  end

  def analyze(store : Store, event : Scanned) : Array(Effect)
    store.scan_counter -= 1
    effects = [] of Effect
    # TODO: should we wait a few ticks to detect files moved during the initial
    # scan
    effects << BeReady.new if store.scan_counter == 0
    effects
  end

  def analyze(store : Store, event : Checksummed) : Array(Effect)
    # Nothing for the moment
    [] of Effect
  end

  def analyze(store : Store, event : FileEvent) : Array(Effect)
    effects = [] of Effect
    if event[:type] == File::Type::Directory
      store.scan_counter += 1
      effects << Scan.new(path: event[:path])
    end
    effects
  end
end
