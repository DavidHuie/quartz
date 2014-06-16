class Quartz::GoProcess

  def initialize(opts)
    @socket_name = "/tmp/quartz#{rand(10000)}.sock"
    ENV['QUARTZ_SOCKET'] = @socket_name

    if opts[:file_path]
      @go_process = Thread.new { `go run #{opts[:file_path]} }` }
    elsif opts[:bin_path]
      @go_process = Thread.new { `#{opts[:bin_path]} }` }
    else
      raise 'Missing go binary path'
    end

    sleep(1)

    @socket = UNIXSocket.new(@socket_name)
  end

  def get_metadata
    payload = {
      'method' => 'Quartz.GetMetadata',
      'params' => [],
      'id' => 1
    }

    @socket.send(payload.to_json, 0)
    read
  end

  def call(struct_name, method, args)
    payload = {
      'method' => 'Quartz.Call',
      'params' => [{
          'StructName' => struct_name,
          'Method' => method,
          'MethodArgs' => args.to_json
        }],
      'id' => 1
    }

    @socket.send(payload.to_json, 0)
    read
  end

  MAX_MESSAGE_SIZE = 1_000_000_000 # Bytes

  def read
    JSON(@socket.recv(MAX_MESSAGE_SIZE))['result']
  end

end
