require "./effect"
require "./event"
require "./store"

module Local
  extend self

  Root = FilePath.new("/")

  def analyze(store : Store, event : TemporalEvent) : Array(Effect)
    effects = [] of Effect
    case event
    when TemporalEvent::Start
      effects << Scan.new(Root)
    when TemporalEvent::Stop
      effects << Close.new
    when TemporalEvent::Tick
      # Nothing for the moment
    else
      raise "Unknown temporal event type: #{event}"
    end
    effects
  end

  def analyze(store : Store, event : FileEvent) : Array(Effect)
    effects = [] of Effect
    if event.type == File::Type::Directory
      effects << Scan.new(event.path)
    end
    effects
  end
end
