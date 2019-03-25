# Working with User Account
`user` function is able to interact with all of your account related settings. 
## Register User Account
First you need to create user account with user name, email and password.
```
$ ankrctl user register user_name --email=user_name@mailinator.com --password=passw0rd 

User user_name@mailinator.com Register Requested, Please Check Your Email Box.
```
Once you registered you should check the email box and confirm the registration with the confirmation code given in the email:

```
$ ankrctl user confirm-registration user_name@mailinator.com --register-code eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTMyNDI0MTgsImp0aSI6IjEzNGZjZjg4LTg4NTUtNGM3Zi05NmZmL
WU1OTU2ZDg4OTZkMCIsImlzcyI6ImFua3IubmV0d29yayJ9.n6InDb5RhwOduTc9-Vo-1rS6CVhqYAF2AnwEnn2B0LE

Confirm Registration Success.
```
## Login User Account
Login to user account before using other function such as `compute` and `wallet`, you should provide user email and password when prompted.

```
$ ankrctl user login

Email: user_name@mailinator.com

Password: passw0rd

Login Successful!
```
## Update Your User Account
You can update some user account properties, such as user name.
```
$ ankrctl user update user_name@mailinator.com --update-key=name --update-value=new_name

User Update Attribute Success.
```

If you have new wallet public key(address), you can update user wallet pubkey with new one, remembet to use the correct wallet public key (address) which is also use in the `wallet` function.
```
$ ankrctl user update user_name@mailinator.com --update-key=pubkey --update-value=<your_new_wallet_pubkey>

User Update Attribute Success.
```
You can change the user password with new one:

```
$ ankrctl user change-password --old-password passw0rd --new-password passw1rd

Change Password Success.
```

You can change user email with new one.
```
$ ankrctl user email-change user_name1@mailinator.com

Change Email Requested, Please Check Your Email Box.
```
Once you requested to change the email you should check the new email box and confirm the change with the confirmation code.

```
$ ankrctl user email-confirm user_name1@mailinator.com --email-code=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ1c2VyX25hbWUxQG1haWxpbmF0b3IuY29tIiwiaXNzIjoiYW5rci5uZXR3b3JrIn
0.DCJAxhryv_kkLOmJLmTWEQrvH1WVOlER0HEJW9CVSj4

Email Change Confirm Success.
```

If you forgot password, you can submit the request to reset:

```
$ ankrctl user forgot-password user_name1@mailinator.com

Forgot Password Requested, Please Check Your Email Box.
```
Once you requested forgot password you should check the email box and confirm the reset with the confirmation code:

```
$ ankrctl user confirm-password user_name1@mailinator.com --password-code=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJ1c2VyX25hbWUxQG1haWxpbmF0b3IuY29tIiwiaXNzIjoiYW5rci5uZXR3
b3JrIn0.DCJAxhryv_kkLOmJLmTWEQrvH1WVOlER0HEJW9CVSj4 --confirm-password=passw0rd

Confirm Password Success.
```
## Logout User Account

```
$ ankrctl user logout

Logout Success.
```