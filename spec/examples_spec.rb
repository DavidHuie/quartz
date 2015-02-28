require 'spec_helper'

describe 'examples' do

  describe 'json' do
    it 'should handle long strings of json' do
      client = Quartz::Client.new(file_path: 'spec/examples/long_json.go')
      res = client[:example].call('LongJson', {})['Output']
      expect(res.count).to eq(8192)
      expect(res.first).to eq("Test")
    end

    it 'should handle many calls' do
      results = []
      client = Quartz::Client.new(file_path: 'spec/examples/long_json.go')

      100.times do
        results.concat(client[:example].call('LongJson', {})['Output'])
      end

      expect(results.count).to eq(819_200)
    end
  end

end
