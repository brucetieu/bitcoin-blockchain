# Deconstructing each piece that makes up a bitcoin address

# pubKey: 084c20702ae5f154334d444a9353679805d9e1422222cf152157992c44e50436dafcc233296f47efee7601f3dd374f345a1af54b458b852527035a8121624871
# pubKeyHash = ripemd160(sha256(pubKey))
b5c9a4a5b98b4948539f728bacf3c2e49c1e2c88
# checksum = sha256(sha256(pubKeyHash))
539fc0b1
# version + pubKeyHash + checksum
00 + b5c9a4a5b98b4948539f728bacf3c2e49c1e2c88 + 539fc0b1 = 00b5c9a4a5b98b4948539f728bacf3c2e49c1e2c88539fc0b1
# address = base58(version + pubKeyHash + checksum)
1HaCs8jFKY7usCgsVhCcAJamEffuVpHgUg