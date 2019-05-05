# Working with Wallet
The Wallet is used to interact with the Ankr's blockchain. Everything that adds new data to the blockchain requires a signature from a private key. These signed transactions can be produced by the `wallet` function.
## Generate Public Key, Private Key and Wallet Address
Wallet address and public key are used for accessing the account tokens on Ankr chain, and private is used to sign the wallet transactions. To use wallet fuction you must create the wallet key pair and save it for later usage.
```
$ ankrctl wallet genkey

Warning: About to generate wallet address, public key and private key..
	Please record and backup wallet address and keys once generated!!
	Note: If these keys lost, you will lost access to your tokens in the wallet!!
	Note: If you have previously generated these keys, the former ones will be replaced!!

	Type 'yes' to confirm that you understood the result of this action: y


Generating keys...

Updating wallet...

Private Key:  WPwpd22hvCz9Erj1GdtpORx18lzZbf4s016+hkVzDBkQxauRZag2XDRQEsziZjOHQpYFXuIQd0q7IJCuiGEzkA==
Public Key:  EMWrkWWoNlw0UBLM4mYzh0KWBV7iEHdKuyCQrohhM5A=
Address:  2966BD78A9E388360BC70B4B5298D08C595D938F
```

## Export Wallet Key Pair to Key File
To better preserve the wallet key pair, we strongly recommend you to export it to key file, and preserve it somewhere safe. If you run ankrctl with docker, remember to copy it to local filesystem using `docker cp` or use mount local filesystem volume option `-v` when using `docker run` to run ankrctl.

When exporting the key file, you need to provide the passcode for key file encryption, passcode should not exceeding 16 characters.
```
$ ankrctl wallet exportkey ./key_file

Warning: About to export privateKey/publicKey/address to key file.
	Type 'yes' to confirm that you would save this key file: yes
Please input the key file encryption secret:
Please input passcode again to confirm:

Exporting keys...

Done.
```
## Import Wallet Key Pair from Key File
After key file exported, you can also import it to ankrctl anytime when you need to access the wallet.

When importing key pair you need to provide the passcode to decrypt the key file.
```
$ ankrctl wallet importkey ./key_file

Warning: About to import address, public key and private key from key file.
	Note: If you have previously generated or imported these keys, the former ones will be replaced!
	Type 'yes' to confirm that you understood the result of this action: yes

Please input the keyfile secret:
Importing...
Private Key:  WPwpd22hvCz9Erj1GdtpORx18lzZbf4s016+hkVzDBkQxauRZag2XDRQEsziZjOHQpYFXuIQd0q7IJCuiGEzkA==
Public Key:  EMWrkWWoNlw0UBLM4mYzh0KWBV7iEHdKuyCQrohhM5A=
Address:  2966BD78A9E388360BC70B4B5298D08C595D938F

Updating wallet...

Done.
```

## Getting Wallet Balance
After you deposit or someone transfer the token to your account, you can query the account balance.
```
$  ankrctl wallet getbal 0D1A90135B1F327FC34BC6515B401A6B19B79125

Query balance by address 0D1A90135B1F327FC34BC6515B401A6B19B79125
The balance is: 6566.123400000000000000
```
## Getting Wallet Balance
If you have tokens at your wallet address and you want to sent the tokens to another account, you can use `sendtoken` and provide the key pair to sign the transaction, valid token format should have no more than 18 digits after decimal point, and not exceeding the balance of your account.
```
$ ankrctl wallet sendtoken 123.56789 --address B508ED0D54597D516A680E7951F18CAD24C7EC9F --target 0D1A90135B1F327FC34BC6515B401A6B19B79125 --public-key=wvHG3EddBbXQHcyJal0CS/YQcNYt
EbFYxejnqf9OhM4= --private-key=wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==

Warning: About to send 123.56789 from wallet address B508ED0D54597D516A680E7951F18CAD24C7EC9F to wallet address 0D1A90135B1F327FC34BC6515B401A6B19B79125, Type 'yes' to confirm this action: yes

Done.
```

## Generate Wallet Address for deposit between MAINNET/ERC20/BEP2
To use Wallet Address for deposit between MAINNET/ERC20/BEP2, use type and purpose to specify these address to generate:
```
go run cmd/ankrctl/main.go -u client-dev.dccn.ankr.com wallet genaddr --type BEP2 --purpose ERC20
Generated Address type BEP2 tbnb15sssy7680ac4726txpzgpzg5tl0v7hh5cxyafj for Purpose ERC20
```
