require 'ruby-prof'
require 'simplecov'

ENV['COVERAGE'] && SimpleCov.start

$LOAD_PATH.unshift(File.join(File.dirname(__FILE__), '..', 'lib'))
$LOAD_PATH.unshift(File.dirname(__FILE__))

require 'rspec'
require 'quartz'

def profile(wall_time = false, &block)
  if wall_time
    RubyProf.measure_mode = RubyProf::WALL_TIME
  end
  RubyProf.start
  call_result = block.call
  result = RubyProf.stop
  File.open('profile', 'w') do |f|
    printer = RubyProf::FlatPrinter.new(result)
    printer.print(f)
  end
  return call_result
end
