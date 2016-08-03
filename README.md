# iam_ssh_authorizedkeys

This is my first `Go` program, and largely WIP and of course pull requests are very much welcome!

The purpose of this project is to provide `AuthorizedKeysCommand` IAM provider for `OpenSSH`. `iam_ssh_authorizedkeys` has been written with [CoreOS](https://coreos.com/) in mind, as it is quite a lot of hassle to install `awscli` on CoreOS in the way that satisfied my bizarre minimalistic taste. 

`iam_ssh_authorizedkeys` does not create users, and assumes that user accounts are created by some other means.

It is expected that your users are already created in `IAM`, and that they have corresponding `SSH keys for AWS CodeCommit` setup. `iam_ssh_authorizedkeys` will iterate through these keys and validate access in real time (as the connection happens).

It is probably worth noting that your `EC2` instances need to be associated with a `IAM Role` that allows the following privileges:
```
iam:ListSSHPublicKeys
iam:GetSSHPublicKey
```

## Usage

iam_ssh_authorizedkeys needs to be placed on the system in location of your chocie (i.e. `/usr/local/bin`)

Also the following lines in `/etc/ssh/sshd_config` are needed:

```
AuthorizedKeysCommand /usr/local/bin/iam_ssh_authorizedkeys
AuthorizedKeysCommandUser nobody
```

Now `ssh user@ec-instance` with a valid key.

## CoreOS cloud-config

For mass deployment use something like this

```
#cloud-config

write_files:
  - path: /etc/ssh/sshd_config
    permissions: 0600
    owner: root:root
    content: |
      # Use most defaults for sshd configuration.
      UsePrivilegeSeparation sandbox
      Subsystem sftp internal-sftp
      ClientAliveInterval 180
      UseDNS no

      PermitRootLogin no
      PasswordAuthentication no
	  AuthorizedKeysCommand /usr/local/bin/iam_ssh_authorizedkeys
	  AuthorizedKeysCommandUser nobody
```

Plus run the following to get the binary on the system

```
curl -fsSL https://github.com/bytewareio/iam_ssh_authorizedkeys/releases/download/iam_authorized_keys-0.0,1/iam_ssh_authorizedkeys-linux64 -o /usr/local/bin/iam_ssh_authorizedkeys
chmod +x /usr/local/bin/iam_ssh_authorizedkeys
```

## Copyright and License

Copyright (c) 2016 BYTEWARE OU <https://www.byteware.io>.
Licensed under the MIT license. See the LICENSE file for full details.

