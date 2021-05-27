from Crypto.PublicKey import RSA
from Crypto.Cipher import PKCS1_OAEP
import binascii

msg = b'{\"payerinn\":\"3528000xxx\",\"payerkpp\":\"997550xxx\",\"payeeinn\":\"3528000xxx\",\"payeekpp\":\"997550xxx\",\"payeebik\":\"044525xxx\",\"payeecheckccount\":\"30101810400000000xxx\",\"payeecorraccount\":\"40702810500020106xxx\",\"amount\":\"1000\",\"details\":\"details\"}'

pubKey = RSA.importKey(open('public.pem').read())
encryptor = PKCS1_OAEP.new(pubKey)
encrypted = encryptor.encrypt(msg)
print("Encrypted:", binascii.hexlify(encrypted))

privKey = RSA.importKey(open('private.pem').read())
decryptor = PKCS1_OAEP.new(privKey)
decrypted = decryptor.decrypt(encrypted)
print('Decrypted:', decrypted)