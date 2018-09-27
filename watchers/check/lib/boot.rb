require 'awesome_print'
require 'fileutils'
require 'pry'

base = File.expand_path "../..", __dir__

Pry.config.history.file = "#{base}/tmp/.pry_history"

FileUtils.cd base do
  FileUtils.mkdir_p "tmp/"
  AwesomePrint.pry!
  require_relative "helpers.rb"
  require_relative "scenario.rb"
  Helpers.cleanup
end
