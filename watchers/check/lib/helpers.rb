require 'fileutils'

module Helpers
  @dir = "tmp"

  class <<self
    attr_reader :dir

    def cleanup
      FileUtils.rm_r "#{@dir}/system-tmp-cozy-drive"
      FileUtils.rm_r "#{@dir}/workspace"
    end
  end
end
