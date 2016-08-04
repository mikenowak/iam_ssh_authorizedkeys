# iam_ssh_authorizedkeys

This is my first `Go` program, and largely WIP and of course pull requests are very much welcome!

The purpose of this project is to provide `AuthorizedKeysCommand` IAM provider for `OpenSSH`. `iam_ssh_authorizedkeys` has been written with [CoreOS](https://coreos.com/) in mind, as it is quite a lot of hassle to install `awscli` on CoreOS in the way that satisfied my bizarre minimalistic taste. 

`iam_ssh_authorizedkeys` does not create users, and assumes that user accounts are created by some other means.

It is expected that your users are already created in `IAM`, and that they have corresponding `SSH keys for AWS CodeCommit` setup. `iam_ssh_authorizedkeys` will iterate through these keys and validate access in real time (as the connection happens).

In order to make this work you need to either:
* Add your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` to `/etc/environment` (less secure)
* Associate your `EC2` instance with an `IAM Role` that has a policy similar to to the bellow attached (more seucre)
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

## Usage

iam_ssh_authorizedkeys needs to be placed on the system in location of your chocie (i.e. `/usr/local/bin`)

Also the following lines in `/etc/ssh/sshd_config` are needed:

```
AuthorizedKeysCommand /usr/local/bin/iam_ssh_authorizedkeys
AuthorizedKeysCommandUser nobody
```

Now `ssh user@ec-instance` with a valid key.

## cloud-config

In order to automate deployment of iam_ssh_authorizedkeys to CoreOS the following `cloud-config` template can be used:

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

## Todo

The following features are currently:
* Support for assuming AWS roles (for cross-account access)

## Copyright and License

Copyright (c) 2016 BYTEWARE OU <https://www.byteware.io>.
Licensed under the MIT license. See the LICENSE file for full details.

