# Working with Wallet
The Wallet is used to interact with the Ankr's blockchain. Everything that adds new data to the blockchain requires a signature from a private key. These signed transactions can be produced by the `wallet` function.
## Generate Public Key, Private Key and Wallet Address
Wallet address and private key are used for accessing the account tokens on Ankr chain, and private key is used to sign the wallet transactions. To use wallet fuction you must create the wallet keystore.
```
$ ankrctl wallet genkey my_new_key

Warning: please record and backup keystore once it is generated, we don’t store your private key!
	 type 'yes' to confirm that you understood the result of this action: y

generating keys...

private key:  AVu+OM3GT3MqISot6GwNRzHb7mdCOSAssCfL5sugtk5Se6mvNvdvKXf96ltbL0rdWTrxfIImbLrFPYr1/ebOIQ==
public key:  Unuprzb3byl3/epbWy9K3Vk68XyCJmy6xT2K9f3mziE=
address:  229FF040112FC1A83D01AA0A43660482C35F6CDF6864F8

about to export to keystore..
please input the keystore encryption password:
please input password again:

exporting to keystore...

created keystore: 
User/my_user/.ankr/UTC--2019-07-24T18-16-12.112674000Z--229FF040112FC1A83D01AA0A43660482C35F6CDF6864F8
```

## Import Wallet Keystore from Keystore file
You can also import keystore file to ankrctl key list anytime, and choose to update keystore address to ankr user account while importing key.
```
$ ankrctl wallet importkey my_import_key --keystore UTC--2019-07-06T22-29-42.344574000Z--219B0A5F896B7A1949128B8F5136362BF939994D769D58

keystore imported: /Users/my_user/.ankr/UTC--2019-07-24T18-34-54.113932000Z--219B0A5F896B7A1949128B8F5136362BF939994D769D58

Warning: do you want to update keystore address of your ankr wallet?
	 type 'yes' to confirm, 'no' to cancel: y

updated user test12345@mailinator.com wallet address: 219B0A5F896B7A1949128B8F5136362BF939994D769D58
```

## List Wallet Keystore
You can get address and public key from keystore list.
```
$ ankrctl wallet keylist
Name             Address                                           Public Key
testkey          AF9BAD386D0F3A78689AC49563F248ECC3A22E956C65B8    B+ajlpUMxgaKPiOD7aANbPH2wno5JifFLKJUhpHZNvg=
testkey1         A0B4B94FF2453DD14402D7856B76CEA8BBCDA3A868A632    pkHtcKIOOkKG0GVl3mpDAsv3bbFdrxxnDhhzHVTSi1k=
```

## Getting Wallet Balance
After you deposit or someone transfer the token to your account, you can query the account balance.
```
$  ankrctl wallet getbalance 0D1A90135B1F327FC34BC6515B401A6B19B79125

Query balance by address 0D1A90135B1F327FC34BC6515B401A6B19B79125
The balance is: 6566.123400000000000000
```
## Send coins
If you have coins at your wallet address and you want to sent the coins to another account, you can use `sendcoins` and provide the keystore to sign the transaction, valid ammount format should have no more than 18 digits after decimal point, and not exceeding the balance of your account.
```
$ ankrctl wallet sendcoins 1.06745756242365 --target-address FB1B2B9561FF55C12FA099C6AF365FE0C88E44D1AC2BCE --keystore my_new_key
Please input the keystore password:

Warning: About to send 1.06745756242365 tokens from address 'A22AF48DA84A984F4EB155F6D811C6331721148B89076B' to address 'FB1B2B9561FF55C12FA099C6AF365FE0C88E44D1AC2BCE', type 'yes' to confirm this action: y

Done.
```

## Generate Wallet Address for deposit between MAINNET/ERC20/BEP2
To use Wallet Address for deposit between MAINNET/ERC20/BEP2, use type and purpose to specify these address to generate:
```
ankrctl wallet genaddr --type BEP2 --purpose ERC20
Generated Address type BEP2 tbnb15sssy7680ac4726txpzgpzg5tl0v7hh5cxyafj for Purpose ERC20
```
