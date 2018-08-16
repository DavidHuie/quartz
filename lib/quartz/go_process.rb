require 'multi_json'

class Quartz::GoProcess

  attr_reader :seed, :socket_path, :temp_file_path

  def self.processes
    @processes ||= []
  end

  def self.clear_processes
    @processes = []
  end

  def forked_mode!
    if @forked.nil?
      @forked = true
      return
    end

    @forked = !@forked
  end

  def initialize(opts)
    @seed = SecureRandom.hex
    socket_dir = opts.fetch(:socket_dir) { '/tmp' }
    @socket_path = File.join(socket_dir, "quartz_#{seed}.sock")
    ENV['QUARTZ_SOCKET'] = @socket_path

    if opts[:file_path]
      Quartz::Validations.check_for_go
      compile_and_run(opts[:file_path])
    elsif opts[:bin_path]
      @go_process = IO.popen(opts[:bin_path])
    else
      raise Quartz::ConfigError, 'Missing go binary'
    end

    block_until_server_starts
    self.class.processes << self
  end

  def compile_and_run(path)
    @temp_file_path = "/tmp/quartz_runner_#{seed}"

    unless system('go', 'build', '-o', @temp_file_path, path)
      raise Quartz::ConfigError, 'Go compilation failed'
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
    max_retries = 20
    retries = 0
    delay = 0.001 # seconds

    loop do
      return if File.exists?(@socket_path)
      raise Quartz::GoServerError, 'RPC server not starting' if retries > max_retries
      sleep(delay * 2**retries)
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

    socket.send(MultiJson.dump(payload), 0)
    response = read

    if response['error']
      raise Quartz::GoResponseError, "Metadata error: #{response['error']}"
    end

    response['result']
  end

  def call(struct_name, method, args)
    payload = {
      'method' => "#{struct_name}.#{method}",
      'params' => [args],
      'id' => 1
    }
    socket.send(MultiJson.dump(payload), 0)
    read
  end

  if ['1.9.3', '2.0.0'].include?(RUBY_VERSION)
    READ_EXCEPTION  = IO::WaitReadable
  else
    READ_EXCEPTION  = IO::EAGAINWaitReadable
  end

  MAX_MESSAGE_SIZE = 8192 # Bytes

  def read
    value = ''
    loop do
      begin
        value << socket.recv_nonblock(MAX_MESSAGE_SIZE)
        break if value.end_with?("\n")
      rescue READ_EXCEPTION
        IO.select([socket], [], [])
      end
    end

    MultiJson.load(value)
  end

  def cleanup
    # If we've forked, there's no need to cleanup since the parent
    # process will.
    return if @forked

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
    return unless @processes
    @processes.each(&:cleanup)
    @processes = []
  end

end

at_exit { Quartz::GoProcess.cleanup }
