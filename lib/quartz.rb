require 'json'
require 'socket'

s = UNIXSocket.new("/tmp/quartz.sock")

payload = {
  'method' => 'Quartz.GetMetadata',
  'params' => [],
  'id' => 1
}

s.send(payload.to_json, 0)
puts s.recv(1000)

payload = {
  'method' => 'Quartz.Call',
  'params' => [{ 'StructName' => 'adder',
                 'Method' => 'Add',
                 'MethodArgs' => { 'A' => 1, 'B' => 2 }.to_json }],
  'id' => 1
}

s.send(payload.to_json, 0)
puts s.recv(1000)
