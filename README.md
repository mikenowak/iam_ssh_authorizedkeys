# iam_ssh_authorizedkeys

This is my first `Go` program, and largely WIP and of course pull requests are very much welcome!

The purpose of this project is to provide `AuthorizedKeysCommand` IAM provider for `OpenSSH`. 

`iam_ssh_authorizedkeys` has been written with [CoreOS](https://coreos.com/) in mind, as it is quite a lot of hassle to install `awscli` on CoreOS in the way that satisfied my bizarre minimalistic taste. 

`iam_ssh_authorizedkeys` does not create system users, and assumes that local user accounts are created by some other means (i.e. `cloud-config`).

## How does it work?

You create your users in `IAM` and upload corresponding `SSH public keys` for these users to IAM under the `Security Credentials` tab.

Then you install 'iam_ssh_authorizedkeys' on your `EC2 instance`, and let it do the magic.

When a user tries to login via ssh, the `sshd` calls `/opt/bin/iam_ssh_authorizedkeys USER` which in its turn makes request to IAM for SSH public keys for a given user. User is then allowed or denied access, based on whether a private key matching one of the public keys returned by `iam_ssh_authorizedkeys` is presented.

## Step by step setup

The following steps relate to [CoreOS](https://coreos.com/), but they really should work on most Linux distributions.

1. Create you users in IAM
2. Upload SSH public key(s) to IAM (this is done in user options under the Security Credentials)
3. Make sure that `ssh_iam_authorizedkeys` has access to valid AWS credentials, one of the following should do the trick:
  * Add your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` to `/etc/environment` (less secure, static credentials)
  * Create an `IAM Role` that has the following policy attached, you use will need to specify this as IAM role when launching your `EC2 instances` (more secure, benefits from auto-rotating credentials)
  ```
{
    "Version": "2012-10-17", 
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListSSHPublicKeys",
                "iam:GetSSHPublicKey"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
  ```

4. Use the following `cloud-config` template as your instance `user-data`

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
      AuthorizedKeysCommand /opt/bin/iam_ssh_authorizedkeys
      AuthorizedKeysCommandUser nobody

coreos:
  units:
    - name: "iamssh.service"
      command: "start"
      content: |
        [Unit]
        Description=Installs iam_ssh_authorizedkeys
        ConditionPathExists=!/opt/bin/iam_ssh_authorizedkeys

        [Service]
        Environment=IAMSSH_PATH=/opt/bin/iam_ssh_authorizedkeys
        Environment=IAMSSH_VER=0.1.0
        Environment=IAMSSH_URL=https://github.com/bytewareio/iam_ssh_authorizedkeys/releases/download/${IAMSSH_VER}/iam_ssh_authorizedkeys-linux64
        Type=oneshot
        RemainAfterExit=yes
        ExecStartPre=/usr/bin/bash -c "mkdir -p /opt/bin && chmod -R 0755 /opt"
        ExecStart=/usr/bin/bash -c "/usr/bin/curl -fsSL --retry 5 --retry-delay 2 -o ${IAMSSH_PATH} ${IAMSSH_URL} && chmod 0755 ${IAMSSH_PATH}"
  ```

5. Launch the instance.

## Caveats

It is important to note that with the setup described above, the users  can add their public keys to their `.ssh/authorized_keys` manually, and they would be still successfully authenticated based on these keys, even if the corresponding public keys don’t exist in IAM.

In order to stop that, and to make sure that only keys added to 'IAM' are consulted for authentication, add the following line to `/etc/ssh/sshd_config`

```
AuthorizedKeysFile /dev/null
```


## Todo

The following features are currently being investigated:
* Support for assuming AWS roles (for cross-account access)

## Copyright and License

Copyright (c) 2016 BYTEWARE OU <https://www.byteware.io>.
Licensed under the MIT license. See the LICENSE file for full details.

