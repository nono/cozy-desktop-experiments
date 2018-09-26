#!/usr/bin/env ruby

require_relative 'lib/boot'
require 'minitest/autorun'
require 'pry-rescue/minitest' unless ENV['CI']

describe 'Testing the local watcher' do
  it 'does some things'
end
