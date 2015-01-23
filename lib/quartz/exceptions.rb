module Quartz
  class Exception < StandardError
  end

  class BadConfig < Exception
  end

  class ArgumentError < Exception
  end
end
