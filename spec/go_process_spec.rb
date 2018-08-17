require 'spec_helper'

describe Quartz::GoProcess do

  let(:process) { Quartz::GoProcess.new(file_path: 'spec/test.go') }

  describe '#get_metadata' do

    context 'with a go file' do
      it 'pulls metadata' do
        expect(process.get_metadata).to eq("adder" => {"NameToMethodMetadata"=>{"Add"=>{"ArgumentToType"=>{"A"=>"int", "B"=>"int"}}, "AddError"=>{"ArgumentToType"=>{"A"=>"int", "B"=>"int"}}}})
      end

      it 'pulls metadata from the recycled socket' do
        new_process = Quartz::GoProcess.new(socket_path: process.socket_path)
        expect(new_process.socket_path).to eq(process.socket_path)
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

  context 'when custom socket directory is used' do

    before { Dir.mkdir(socket_dir) unless File.directory?(socket_dir) }

    let(:socket_dir) { '/tmp/x' }
    let(:process) { Quartz::GoProcess.new(file_path: 'spec/test.go', socket_dir: socket_dir) }

    it 'works with custom socket directory' do
      expect(File.exists?(process.socket_path)).to be_truthy
      process.cleanup
      expect(File.exists?(process.socket_path)).to be_falsey
    end

    it 'creates the socket in the socket dir' do
      expect(process.socket_path).to match(/^#{socket_dir}\/quartz_[a-f\d]+\.sock/)
    end

  end

  context 'with a custom socket path' do

    let(:new_process) { Quartz::GoProcess.new(socket_path: process.socket_path) }

    it 'the new process does not clean up the existing socket' do
      expect(File.exists?(new_process.socket_path)).to be_truthy
      new_process.cleanup
      expect(File.exists?(new_process.socket_path)).to be_truthy
    end

    it 'creates the socket in the socket dir' do
      result = new_process.call('adder', 'Add', { 'A' => 5, 'B' => 6 })
      expect(result).to eq({"id"=>1, "result"=>{"X"=>11}, "error"=>nil})
    end

  end

  describe '#forked_mode' do
    it 'creates a new socket' do
      socket = process.socket
      process.forked_mode!
      expect(process.socket).not_to eq(socket)
    end
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
        expect(File.exists?(process.temp_file_path)).to be_falsey
        expect($?.exited?).to be_truthy
      end

      it 'does not kill child processes if forked' do
        expect(File.exists?(process.temp_file_path)).to be_truthy
        process
        process.forked_mode!
        process.cleanup
        expect(File.exists?(process.temp_file_path)).to be_truthy

        process.forked_mode!
        process.cleanup
        expect(File.exists?(process.temp_file_path)).to be_falsey
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
