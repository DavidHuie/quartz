class Quartz::GoProcess

  attr_reader :temp_file_path

  def self.processes
    @processes ||= []
  end

  def seed
    @seed
  end

  def initialize(opts)
    @seed = rand(10000)
    @socket_path = "/tmp/quartz_#{seed}.sock"
    ENV['QUARTZ_SOCKET'] = @socket_path

    if opts[:file_path]
      compile_and_run(opts[:file_path])
    elsif opts[:bin_path]
      @go_process = IO.popen(opts[:bin_path])
    else
      raise 'Missing go binary'
    end

    block_until_server_starts
    self.class.processes << self
  end

  def compile_and_run(path)
    @temp_file_path = "/tmp/quartz_runner_#{seed}"

    unless system('go', 'build', '-o', @temp_file_path, path)
      raise 'Go compilation failed'
    end

    @go_process = IO.popen(@temp_file_path)
  end

  def pid
    @go_process.pid
  end

  def socket
    Thread.current["quartz_socket_#{seed}".to_sym] ||= UNIXSocket.new(@socket_path)
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
      # This parameter isn't needed because we use a different
      # connection for each thread.
      'id' => 1
    }

    socket.send(payload.to_json, 0)
    response = read

    if response['error']
      raise "Metadata error: #{read['error']}"
    end

    response['result']
  end

  def call(struct_name, method, args)
    payload = {
      'method' => "#{struct_name}.#{method}",
      'params' => [args],
      'id' => 1
    }
    socket.send(payload.to_json, 0)
    read
  end

  MAX_MESSAGE_SIZE = 1_000_000_000 # Bytes

  def read
    JSON(socket.recv(MAX_MESSAGE_SIZE))
  end

  def cleanup
    unless @killed_go_process
      Process.kill('SIGTERM', pid)
      Process.wait(pid)
      @killed_go_process = true
    end

    if @temp_file_path && File.exists?(@temp_file_path)
      File.delete(@temp_file_path)
    end
  end

  def self.cleanup
    @processes.each(&:cleanup)
    @processes = []
  end

end

at_exit { Quartz::GoProcess.cleanup }
