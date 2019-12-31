require "inotify"

# TODO: Write documentation for `Watcher`
module Watcher
  def self.start(dir)
    watcher = Inotify.watch dir do |event|
      pp! event
    end

    # ... for 10 seconds.
    sleep 10.seconds
    watcher.close
  end
end
