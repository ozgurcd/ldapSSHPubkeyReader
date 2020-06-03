# ldapSSHPubkeyReader

OpenSSH relies on an external application to provide SSH keys when AuthorizedKeysCommand directive used. In order to use SSH with keys stored in LDAP, a suitable script or program needs to provide the keys. Traditionally, some shell script that calls ldapsearch is being used for this purpose.

This go program can be used to fullfill the same need, only faster and in a way lightweight manner.

Notes:

This program requires one argument as ldap uid and dumps all of the users public keys to stdout, as sshd requires. Also, currently it is ignoring any SSL based errors (such as unknown CA), I am planning to make it optional in the future.

It reads the configuration file "ldapPubKeyReader.json" from the same directory it is located.
