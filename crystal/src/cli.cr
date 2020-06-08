require "./desktop"
require "./local/recorder"
require "./web/dev"

# TODO: parse option for argv = ARGV[1..]
case ARGV.first?
when "record"
  dir = File.expand_path "../tmp/Cozy", __DIR__
  recorder = Local::Recorder.new dir
  Signal::INT.trap do
    STDERR.puts "Exit"
    recorder.stop
    exit
  end
  recorder.start
when "sync"
  dir = File.expand_path "../tmp/Cozy", __DIR__
  runner = Desktop.create_runner dir
  Signal::INT.trap do
    STDERR.puts "Exit"
    runner.stop
    runner.wait_final_stop
    exit
  end
  runner.run
when "web"
  Signal::INT.trap do
    STDERR.puts "Exit"
    Kemal.stop
    exit
  end
  Web::Dev.run
else
  puts "Usage: #{PROGRAM_NAME} sync"
end
