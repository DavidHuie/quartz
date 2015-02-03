module Quartz
  class Exception < StandardError; end

  class ResponseError < Exception; end
  class ConfigError < Exception; end
  class ArgumentError < Exception; end
  class GoServerError < Exception; end
  class GoResponseError < Exception; end
end
