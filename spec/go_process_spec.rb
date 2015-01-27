require 'spec_helper'

describe Quartz::GoProcess do

  let(:process) { Quartz::GoProcess.new(file_path: 'spec/test.go') }

  describe '#get_metadata' do

    context 'with a go file' do

      it 'pulls metadata' do
        expect(process.get_metadata).to eq("adder" => {"NameToMethodMetadata"=>{"Add"=>{"ArgumentToType"=>{"A"=>"int", "B"=>"int"}}, "AddError"=>{"ArgumentToType"=>{"A"=>"int", "B"=>"int"}}}})
      end

    end

  end

  describe '#call' do

    it 'is able to call a method on a struct' do
      result = process.call('adder', 'Add', { 'A' => 5, 'B' => 6 })
      expect(result).to eq({"id"=>1, "result"=>{"X"=>11}, "error"=>nil})
    end

  end

  it 'works with an existing binary file' do
    temp_file = "/tmp/quartz_test_#{rand(10000)}"
    system('go', 'build', '-o', temp_file, 'spec/test.go')
    process = Quartz::GoProcess.new(bin_path: temp_file)
    result = process.call('adder', 'Add', { 'A' => 5, 'B' => 6 })
    expect(result).to eq({"id"=>1, "result"=>{"X"=>11}, "error"=>nil})
    File.delete(temp_file)
  end

  describe '.cleanup' do

    context 'files' do

      it "it deletes temporary files" do
        expect(File.exists?(process.temp_file_path)).to be_truthy
        process.cleanup
        expect(File.exists?(process.temp_file_path)).to be_falsey
      end

    end

    context 'processes' do

      it "it kills child processes" do
        expect(File.exists?(process.temp_file_path)).to be_truthy
        process.cleanup
        expect($?.exited?).to be_truthy
      end

    end

    context 'sockets' do

      it 'cleans up sockets created by the go application' do
        process.cleanup
        expect(File.exists?(process.socket_path)).to be_falsey
      end

    end

  end

end
