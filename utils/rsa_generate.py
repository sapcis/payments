from Crypto.PublicKey import RSA

keyPair = RSA.generate(4096)

pubKey = keyPair.publickey()
pubKeyPEM = pubKey.exportKey('PEM')
print(pubKeyPEM.decode('ascii'))
f = open('public.pem','wb')
f.write(pubKeyPEM)
f.close()

privKeyPEM = keyPair.exportKey('PEM')
print(privKeyPEM.decode('ascii'))
f = open('private.pem','wb')
f.write(privKeyPEM)
f.close()