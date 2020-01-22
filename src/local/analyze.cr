require "./effect"
require "./event"
require "./store"

module Local
  extend self

  Root = FilePath.new("/")

  def analyze(store : Store, event : TemporalEvent::Start) : Array(Effect)
    store.scan_counter += 1
    [Scan.new(Root)] of Effect
  end

  def analyze(store : Store, event : TemporalEvent::Stop) : Array(Effect)
    [Close.new] of Effect
  end

  def analyze(store : Store, event : TemporalEvent::Tick) : Array(Effect)
    # Nothing for the moment
  end

  def analyze(store : Store, event : TemporalEvent::Scanned) : Array(Effect)
    store.scan_counter -= 1
    effects = [] of Effect
    effects << BeReady.new if store.scan_counter == 0
    effects
  end

  def analyze(store : Store, event : FileEvent) : Array(Effect)
    effects = [] of Effect
    if event.type == File::Type::Directory
      store.scan_counter += 1
      effects << Scan.new(event.path)
    end
    effects
  end
end
