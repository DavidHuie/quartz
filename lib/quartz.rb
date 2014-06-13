require 'json'
require 'socket'

s = UNIXSocket.new("/tmp/quartz.sock")

payload = {
  'method' => 'Quartz.Add',
  'params' => [{ 'N' => 2, 'M' => 3}],
  'id' => 4
}

s.send(payload.to_json, 0)
puts s.recv(100)
