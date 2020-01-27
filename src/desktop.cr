require "./file_info"
require "./local/component"
require "./remote/component"
require "./runner"
require "./sync/component"

# TODO: Write documentation for `Desktop`
module Desktop
  VERSION = "0.1.0"

  def self.start(dir : String)
    local = Local::Component.new dir
    remote = Remote::Component.new
    sync = Sync::Component.new
    runner = Runner.new local: local, remote: remote, sync: sync
    runner.run
    runner
  end
end
