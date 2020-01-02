module Local
  class Analyzer
    struct Response
      property pending, cksums, stats

      def initialize(@pending : UInt32, @cksums : Array(FilePath), @stats : Array(FilePath))
      end
    end

    def initialize(@store : Local::Store)
    end

    def receive(event : Local::Event)
    end
  end
end
