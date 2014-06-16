require 'quartz'

go_process = Quartz::GoProcess.new(file_path: 'lookup_dns.go')
client = Quartz::Client.new(go_process)

puts "Structs: #{client.structs}"
puts "Struct methods for #{client[:resolver].struct_name}: #{client[:resolver].struct_methods}"

puts client[:resolver].call('FindIPs',
                            'Hostnames' => ['www.google.com',
                                            'www.facebook.com',
                                            'www.microsoft.com'])
