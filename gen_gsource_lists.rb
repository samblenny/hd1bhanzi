#!/usr/bin/ruby
require "set"

C20_IN = "Unihan/Unihan_DictionaryLikeData.txt"
IIC_IN = "Unihan/Unihan_IRGSources.txt"
C20_OUT = "core2020_g.txt"
IIC_OUT = "iicore_g.txt"
C20_RE = Regexp.compile(/U\+(.*)\tkUnihanCore2020.*G.*/)
IIC_RE = Regexp.compile(/U\+(.*)\tkIICore.*G.*/)

# Write file of UnihanCore2020 Source G (Simplified Chinese) character hex codepoints
c20 = File.read(C20_IN).lines
        .select { |n| C20_RE.match n }
        .map { |n| C20_RE.match(n)[1] }
puts "Writing #{c20.size} characters to #{C20_OUT}"
File.open(C20_OUT, "w") { |f| f.puts(c20) }

# Write file of IICore Source G (Simplified Chinese) character hex codepoints
iic = File.read(IIC_IN).lines
        .select { |n| IIC_RE.match n }
        .map { |n| IIC_RE.match(n)[1] }
puts "Writing #{iic.size} characters to #{IIC_OUT}"
File.open(IIC_OUT, "w") { |f| f.puts(iic) }

# Compare overlap
iic_set = Set.new(iic)
c20_set = Set.new(c20)
puts "(iic - c20).size = #{(iic_set-c20_set).size}"
puts "(c20 - iic).size = #{(c20_set-iic_set).size}"
