class Quartz::Client

  def initialize(go_process)
    @process = go_process
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
