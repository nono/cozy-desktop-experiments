require 'fileutils'

module Helpers
  @dir = "tmp"

  class <<self
    attr_reader :dir

    def cleanup
      ["#{@dir}/system-tmp-cozy-drive", "#{@dir}/workspace"].each do |d|
        FileUtils.rm_r d if File.exist? d
      end
    end
  end
end
