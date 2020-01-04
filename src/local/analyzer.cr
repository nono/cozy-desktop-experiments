require "./event"
require "./store"

module Local
  class Analyzer
    struct Response
      property pending, cksums, stats

      def initialize(@pending : UInt32, @cksums : Array(FilePath), @stats : Array(FilePath))
      end
    end

    def initialize(@store : Store)
    end

    def receive(event : Event)
    end
  end
end
