require "./operation"

module Simulator
  struct Scenario
    alias Client = NamedTuple(name: String, OS: String)

    property pending : String?
    property clients : Array(Client)
    property ops : Array(Operation)

    def initialize(nb_clients)
      case nb_clients
      when 1
        @clients = [{name: "desktop", os: "linux"}]
      when 2
        @clients = [{name: "desktop", os: "linux"},
                    {name: "laptop", os: "linux"}]
      else
        @clients = [] of String
        nb_clients.times do
          @clients << {name: Faker::Internet.user_name, os: "linux"}
        end
      end
      @ops = [] of Operation
    end
  end
end
