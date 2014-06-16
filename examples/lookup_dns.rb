require 'quartz'

client = Quartz::Client.new(file_path: 'lookup_dns.go')

puts "Structs: #{client.structs}"
puts "Struct methods for #{client[:resolver].struct_name}: #{client[:resolver].struct_methods}"
puts "Response from FindIPs:"
puts client[:resolver].call('FindIPs',
                            'Hostnames' => ['www.google.com',
                                            'www.facebook.com',
                                            'www.microsoft.com'])
