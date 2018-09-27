#!/usr/bin/env ruby

require_relative 'lib/boot'
require 'minitest/autorun'
require 'rantly/minitest_extensions'
require 'pry-rescue/minitest' unless ENV['CI']

describe "Cozy Drive" do
  it "watches correctly the local file-system" do
    property_of { operations }.check(1) do |ops|
      scenario = Scenario.new ops
      scenario.play
      ap scenario.results
      # assert_kind_of Integer, scenario, "integer property did not return Integer type"
    end
  end
end
