require 'spec_helper'

describe 'examples' do

  describe 'json' do
    it 'should handle long strings of json' do
      client = Quartz::Client.new(file_path: 'spec/examples/long_json.go')
      res = client[:example].call('LongJson', {})['Output']
      expect(res.count).to eq(8192)
      expect(res.first).to eq("Test")
    end
  end

end
