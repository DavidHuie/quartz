require 'spec_helper'

describe Quartz::Client do

  let(:client) { Quartz::Client.new(file_path: 'spec/test.go') }

  it 'creates structs internally' do
    result = client[:adder].call('Add', 'A' => 2, 'B' => 5)
    expect(result).to eq({'X' => 7})
  end

end
