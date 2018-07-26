all: secp256k1 adx

adx:
	go install github.com/tokenme/adx

secp256k1:
	cp -r dependencies/secp256k1/src vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1/src;
	cp -r dependencies/secp256k1/include vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1/include;

install:
	rm -rf /opt/adx-ui/*;
	cp -r ui/build/dist/* /opt/adx-ui/;
	rm -rf /opt/adx-templates/*;
	cp -r templates/* /opt/adx-templates/;
	cp -f /opt/go/bin/adx /usr/local/bin/;
	chmod a+x /usr/local/bin/adx;
