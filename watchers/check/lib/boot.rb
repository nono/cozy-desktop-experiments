require 'awesome_print'
require 'fileutils'
require 'pry'

base = File.expand_path "../..", __dir__
FileUtils.cd base do
  FileUtils.mkdir_p "tmp/"
  AwesomePrint.pry!
  Pry.config.history.file = "tmp/.pry_history"
  require_relative "helpers.rb"
  require_relative "model.rb"
  Helpers.cleanup
end
