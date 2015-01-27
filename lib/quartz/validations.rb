module Quartz::Validations

  GO_FILE_LOCATION = 'src/github.com/DavidHuie/quartz/go/quartz/quartz.go'

  def self.check_for_go
    go_exists = ENV['PATH'].split(File::PATH_SEPARATOR).any? do |directory|
      File.exist?(File.join(directory, 'go'))
    end

    raise 'Go not installed.' unless go_exists
  end

  def self.check_go_quartz_version
    current_quartz = File.read(File.join(File.dirname(__FILE__), '../../go/quartz/quartz.go'))

    installed_quartz_dir = ENV['GOPATH'].split(File::PATH_SEPARATOR).select do |directory|
      File.exist?(File.join(directory, Quartz::Validations::GO_FILE_LOCATION))
    end[0]

    unless installed_quartz_dir
      raise "GOPATH not configured."
    end

    installed_quartz = File.read(File.join(installed_quartz_dir,
                                          Quartz::Validations::GO_FILE_LOCATION))

    if current_quartz != installed_quartz
      STDERR.write <<-EOS
Warning: the version of Quartz in $GOPATH does not match
the version packaged with the gem. Please update the Go
Quartz package.
EOS
    end
  end

end
