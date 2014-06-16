require 'spec_helper'

describe Quartz::Client do

  let(:process) { Quartz::GoProcess.new(file_path: 'spec/test.go') }

  it 'creates structs internally' do
    c = Quartz::Client.new process
    result = c[:adder].call('Add', 'A' => 2, 'B' => 5)
    result.should eq({'X' => 7})
  end

end
