#!/usr/bin/ruby
# coding: utf-8
require "set"

C20_IN = "Unihan/Unihan_DictionaryLikeData.txt"
IIC_IN = "Unihan/Unihan_IRGSources.txt"
C20_OUT = "core2020_g.txt"
IIC_OUT = "iicore_g.txt"
C20_RE = Regexp.compile(/U\+(.*)\tkUnihanCore2020.*G.*/)
IIC_RE = Regexp.compile(/U\+(.*)\tkIICore.*G.*/)
PUNCTUATION_IN = "cjk_punctuation.txt"
INDEX_OUT = "hanzi_core2020_g_index.txt"

# This is only read to verify against the other lists; data from this list is not copied
TGH_IN = "china-8105-06062014.txt"
TGH_RE = Regexp.compile(/U\+(.*)\t[0-9]+/)

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

# Cross check against TGH-2013
tgh = File.read(TGH_IN).lines
        .select { |n| TGH_RE.match n }
        .map { |n| TGH_RE.match(n)[1] }
puts "TGH list has #{tgh.size} characters"
tgh_set = Set.new(tgh)
puts "(tgh - iic).size = #{(tgh_set-iic_set).size}"
puts "(iic - tgh).size = #{(iic_set-tgh_set).size}"
puts "(c20 - tgh).size = #{(c20_set-tgh_set).size}"
puts "(tgh - c20).size = #{(tgh_set-c20_set).size}"

puts "iic - c20:"
puts (iic - c20).map { |n| n.to_i(16).chr('UTF-8') }.join("")
puts "iic - tgh:"
puts (iic - tgh).map { |n| n.to_i(16).chr('UTF-8') }.join("")
puts "c20 - tgh:"
puts (c20 - tgh).map { |n| n.to_i(16).chr('UTF-8') }.join("")

# Create index file with UnihanCore2020 Source G + CJK punctuation
c20 = File.read(C20_OUT).strip
punct = File.read(PUNCTUATION_IN).strip
File.open(INDEX_OUT, "w") { |f|
  f.puts(c20)
  f.puts(punct)
}

# Expected output:
#
# Writing 8249 characters to core2020_g.txt
# Writing 5825 characters to iicore_g.txt
# (iic - c20).size = 22
# (c20 - iic).size = 2446
# TGH list has 8105 characters
# (tgh - iic).size = 2322
# (iic - tgh).size = 42
# (c20 - tgh).size = 144
# (tgh - c20).size = 0
# iic - c20:
# 伕劻卻坵屌晥晳枓洩濛濬珮甽睪矇硃礽穀肏蹠閒鯈
# iic - tgh:
# 䢵䣅䣓䧑伕劻卻囝坵垅嬲屌拚摺晥晳枓洩濛濬珮甙甽睪矇矽砦硃硷磺礽穀肏舨菸蹓蹠閒饤饾馀鯈
# c20 - tgh:
# 䢵䣅䣓䧑丌丨丬丶丿乇亠亻佧冂冖冫凵刂勹匚卩厶咭哜囗囝坶垅堀塃塈塮夂嬲宀屮巛庀廴廾弪彐
# 彡後忄憝扌拚挢捱揎揞揲揾搿摺擗攴攵朊朘楱榀榘檵欷氵氽沲泶渖湨灬熳犭猓猸甙疋疒痖矽砇砦
# 砩硷磺礤礳礻穦窆箝簦糸纟缋缍缏耲肀肜膣膪臁舡舨艹茇菸葓虍蚵螓蟓蠛衤謦讠诶谘跫蹓軎辶醣
# 钅钶钸镙镟阝阢飚饣饤饾馀髟鲶鳆鳋麴齄
