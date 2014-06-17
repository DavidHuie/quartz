class Quartz::GoProcess

  def initialize(opts)
    @socket_path = "/tmp/quartz_#{rand(10000)}.sock"
    ENV['QUARTZ_SOCKET'] = @socket_path

    if opts[:file_path]
      compile_and_run(opts[:file_path])
    else
      raise 'Missing go binary'
    end

    block_until_server_starts
    register_pid
    register_temp_file
  end

  def compile_and_run(path)
    @temp_file_path = "/tmp/quartz_runner_#{rand(10000)}"

    unless system('go', 'build', '-o', @temp_file_path, path)
      raise 'Go compilation failed'
    end

    @go_process = IO.popen(@temp_file_path)
  end

  def self.temp_files
    @temp_files ||= []
  end

  def self.child_pids
    @child_pids ||= []
  end

  def register_pid
    self.class.child_pids << @go_process.pid
  end

  def register_temp_file
    self.class.temp_files << @temp_file_path
  end

  def socket
    Thread.current[:quartz_socket] ||= UNIXSocket.new(@socket_path)
  end

  def block_until_server_starts
    max_retries = 10
    retries = 0
    delay = 0.1 # seconds

    loop do
      raise 'RPC server not starting' if retries > max_retries
      return if File.exists?(@socket_path)
      sleep(delay * retries * 2**retries)
      retries += 1
    end
  end

  def get_metadata
    payload = {
      'method' => 'Quartz.GetMetadata',
      'params' => [],
      'id' => 1
    }

    socket.send(payload.to_json, 0)
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

    socket.send(payload.to_json, 0)
    read
  end

  MAX_MESSAGE_SIZE = 1_000_000_000 # Bytes

  def read
    JSON(socket.recv(MAX_MESSAGE_SIZE))['result']
  end

end

at_exit do
  Quartz::GoProcess.child_pids.each do |pid|
    Process.kill('SIGTERM', pid)
  end

  Process.waitall

  Quartz::GoProcess.temp_files.each do |file_path|
    File.delete(file_path)
  end
end
