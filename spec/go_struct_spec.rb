require 'spec_helper'

describe Quartz::GoStruct do

  let(:process) { Quartz::GoProcess.new(file_path: 'spec/test.go') }

  describe '#call' do

    let(:struct) do
      Quartz::GoStruct.new("adder",
                           {"NameToMethodMetadata"=>{"Add"=>{"ArgumentToType"=>{"A"=>"int", "B"=>"int"}}}}, process)
    end

    it 'is able to call a struct' do
      response = struct.call('Add', 'A' => 2, 'B' => 4)
      response.should eq({'X' => 6})
    end

  end

end
