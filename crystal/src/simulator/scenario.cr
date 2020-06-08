require "json"
require "./operation"

module Simulator
  class Scenario
    alias Client = NamedTuple(name: String, os: String)
    alias CozyInstance = NamedTuple(active: Bool)

    JSON.mapping(
      pending: String?,
      clients: Array(Client),
      cozy: CozyInstance,
      ops: Array(Operation)
    )

    def initialize(nb_clients : Int)
      case nb_clients
      when 0, 1
        @clients = [{name: "desktop", os: "linux"}]
      when 2
        @clients = [{name: "desktop", os: "linux"},
                    {name: "laptop", os: "linux"}]
      else
        @clients = [] of Client
        nb_clients.times do
          @clients << {name: Faker::Internet.user_name, os: "linux"}
        end
      end
      @cozy = {active: nb_clients > 0}
      @ops = [] of Operation
    end
  end
end
