require "./operation"

module Simulator
  struct Scenario
    property pending : String?
    property clients : Array(String)
    property ops : Array(Operation)

    def initialize(nb_clients)
      case nb_clients
      when 1
        @clients = ["desktop"]
      when 2
        @clients = ["desktop", "laptop"]
      else
        @clients = [] of String
        nb_clients.times do
          @clients << Faker::Internet.user_name
        end
      end
      @ops = [] of Operation
    end
  end
end
