module Quartz::Validations

  def self.check_for_go
    go_exists = ENV['PATH'].split(File::PATH_SEPARATOR).any? do |directory|
      File.exist?(File.join(directory, 'go'))
    end

    raise 'Go not installed' unless go_exists
  end

  def self.check_go_quartz_version
    current_pulse = File.read(File.join(File.dirname(__FILE__), '../../go/quartz/quartz.go'))
    installed_pulse = File.read(File.join(ENV['GOPATH'],
                                          'src/github.com/DavidHuie/quartz/go/quartz/quartz.go'))

    if current_pulse != installed_pulse
      STDERR.write <<-EOS
Warning: the version of Quartz in $GOPATH does not match
the version packaged with the gem. Please update the Go
Quartz package.
EOS
    end
  end

end
