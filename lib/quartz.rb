require 'json'
require 'securerandom'
require 'socket'

module Quartz
  class ResponseError < StandardError; end
end

require 'quartz/go_process'
require 'quartz/go_struct'
require 'quartz/client'
require 'quartz/validations'

Quartz::Validations.check_for_go
Quartz::Validations.check_go_quartz_version
