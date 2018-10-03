require 'rantly'

class Rantly
  def operations
    @knows_path = []
    k = range 1, 10
    init = Array.new(k) { init_op }
    n = range 5, 10
    ops = Array.new(n) { op }
    init + [:start_client] + ops
  end

  def init_op
    [choose(:mkdir, :create_file), new_path]
  end

  def op
    freq [3, :sleep],
         [2, :mkdir],
         [2, :create_file],
         [1, :update_file]
  end

  def sleep
    k = range(0, 4)
    [:sleep, 2**k]
  end

  def mkdir
    [:mkdir, path]
  end

  def create_file
    [:create_file, path, content]
  end

  def update_file
    [:update_file, existing_path, content]
  end

  def path
    freq [2, :new_path],
         [1, :existing_path]
  end

  def new_path
    name = string
    guard !(name.include? '/')
    @knows_path << name
    name
  end

  def existing_path
    i = integer % @knows_path.length
    @knows_path[i]
  end

  def content
    sized(range(10, 100)) { string }
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
    @operations.each do |op|
      ap op
    end
    @results = []
  end
end
