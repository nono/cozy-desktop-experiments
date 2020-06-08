require "./spec_helper"
require "../../src/simulator/*"

describe Simulator::Scenario do
  it "can be encoded and parsed as JSON" do
    nb_clients = Random.rand 3
    scenario = Simulator::Generate.new(nb_clients).run
    serialized = scenario.to_json
    parsed = Simulator::Scenario.from_json serialized
    parsed.pending.should eq scenario.pending
    parsed.clients.should eq scenario.clients
    parsed.cozy.should eq scenario.cozy
    parsed.ops.size.should eq scenario.ops.size
    parsed.ops.each_with_index do |op, i|
      op.should eq scenario.ops[i]
    end
  end
end
