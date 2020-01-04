require "./desktop"
require "./local/recorder"

# TODO: parse option for argv = ARGV[1..]
case ARGV.first?
when "record"
  dir = File.expand_path "../tmp/Cozy", __DIR__
  recorder = Local::Recorder.start dir
  Signal::INT.trap do
    STDERR.puts "Exit"
    recorder.stop
  end
when "sync"
  raise "Not yet implemented"
else
  puts "Usage: #{PROGRAM_NAME} sync"
end
