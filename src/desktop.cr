require "inotify"

# TODO: Write documentation for `Desktop`
module Desktop
  VERSION = "0.1.0"

  def self.start(dir)
    watcher = Inotify.watch dir do |event|
      pp! event
    end

    # ... for 10 seconds.
    sleep 10.seconds
    watcher.close
  end
end
