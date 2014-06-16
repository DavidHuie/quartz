require 'quartz'

go_process = Quartz::GoProcess.new(file_path: 'lookup_dns.go')
client = Quartz::Client.new(go_process)

# Go ahead and try to resolve a few million hostnames ;)
puts client[:resolver].call('FindIPs',
                            'Hostnames' => ['www.google.com',
                                            'www.facebook.com',
                                            'www.microsoft.com'])
