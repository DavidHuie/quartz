class Quartz::Client

  def initialize(opts)
    @process = Quartz::GoProcess.new(file_path: opts[:file_path])
    @structs = {}
    @process.get_metadata.each do |struct_name, metadata|
      @structs[struct_name.to_sym] = Quartz::GoStruct.new(struct_name, metadata, @process)
    end
  end

  def [](struct_name)
    @structs[struct_name]
  end

  def structs
    @structs.keys
  end

end
