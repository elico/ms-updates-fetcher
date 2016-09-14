#!/usr/bin/env ruby

require "digest"
require "open-uri"

request = URI.parse(ARGV[1])

puts ARGV[0]
puts ARGV[1]

#http://fg.v4.download.windowsupdate.com/c/msdownload/update/software/crup/2016/06/windows10.0-kb3163019-x64_e58c1a23ae5b85b00d79c146e42dd99f3a949c8b.cab


puts "public URL KEY => "+ ARGV[0]+":"+ARGV[1]
puts Digest::SHA256.hexdigest ARGV[0]+":"+ARGV[1]

if request.host.end_with?("download.windowsupdate.com")
	request.host = "msupdates.ngtech.internal"
end

puts "private URL KEY => "+ ARGV[0]+":"+request.to_s
puts Digest::SHA256.hexdigest ARGV[0]+":"+request.to_s
