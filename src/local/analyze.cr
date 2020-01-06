require "./effect"
require "./event"
require "./store"

module Local
  Root = FilePath.new("/")

  def self.analyze(store : Store, event : Event) : Array(Effect)
    effects = [] of Effect
    case event
    when TemporalEvent::Start
      effects << Scan.new(Root)
    when TemporalEvent::Stop
      effects << Close.new
    end
    effects
  end
end
