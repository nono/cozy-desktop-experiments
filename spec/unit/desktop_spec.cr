require "./spec_helper"

describe Desktop do
  it "has a version number" do
    Desktop::VERSION.should match(/^\d+\.\d+\.\d+$/)
  end
end
