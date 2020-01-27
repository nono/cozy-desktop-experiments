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
  dir = File.expand_path "../tmp/Cozy", __DIR__
  runner = Desktop.start dir
  Signal::INT.trap do
    STDERR.puts "Exit"
    runner.stop
  end
else
  puts "Usage: #{PROGRAM_NAME} sync"
end
