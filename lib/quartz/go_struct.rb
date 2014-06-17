class Quartz::GoStruct

  attr_reader :struct_name, :struct_methods

  def initialize(struct_name, method_info, process)
    @struct_name = struct_name
    @method_name_to_arg_info = {}
    @process = process

    method_info["NameToMethodMetadata"].each do |method_name, info|
      @method_name_to_arg_info[method_name] = info["ArgumentToType"].keys()
    end

    @struct_methods = @method_name_to_arg_info.keys
  end

  def call(method_name, args)
    unless @struct_methods.include?(method_name)
      raise "Invalid method: #{method_name}"
    end

    arg_info = @method_name_to_arg_info[method_name]

    # Validate arguments
    args.each do |k, v|
      unless arg_info.include?(k)
        raise "Invalid argument: #{k}"
      end

      # TODO: validate type
    end

    response = @process.call(@struct_name, method_name, args)

    if response['error']
      raise "Error calling #{method_name}: #{response['error']}"
    end

    response['result']
  end

end
