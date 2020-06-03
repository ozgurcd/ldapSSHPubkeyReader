# ldapSSHPubkeyReader

This go program can be used to retrieve SSH Keys from a LDAP server in order to use in SSH configuration with AuthorizedKeysCommand directive.


Notes:

This program requires one argument as ldap uid and dumps all of the users public keys to stdout, as sshd requires. Also, currently it is ignoring any SSL based errors (such as unknown CA), I am planning to make it optional in the future.
