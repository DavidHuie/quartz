require 'json'
require 'socket'

module Quartz
  class ResponseError < StandardError; end
end

require 'quartz/go_process'
require 'quartz/go_struct'
require 'quartz/client'

# Check if go is installed
go_exists = ENV['PATH'].split(File::PATH_SEPARATOR).any? do |directory|
  File.exist?(File.join(directory, 'go'))
end

raise 'Go not installed' unless go_exists
