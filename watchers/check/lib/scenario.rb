require 'rantly'

class Rantly
  def operations
    integer
  end
end

class Scenario
  attr_accessor :operations, :results

  def self.create
    ops = Rantly { operations }
    new ops
  end

  def initialize(operations)
    @operations = operations
  end

  def play
    @results = []
  end
end
