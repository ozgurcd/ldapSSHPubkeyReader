# ldapSSHPubkeyReader

OpenSSH relies on an external application to provide SSH keys when AuthorizedKeysCommand directive is used. In order to use SSH with keys stored in LDAP, a suitable script or program needs to provide the keys. Traditionally, some shell script that calls ldapsearch is being used for this purpose.

This go program can be used to fullfill the same need, only faster and in a lightweight manner.

### Compile

Please modify the provided Makefile according to your OS and Architecture. Then issue the command "make", if you want to create a statically compile binary, issue the command "make static".

Please note, LDFLAGS is currently set to -ldflags '-w -s', which produces a smaller binary but strips out debug information. If you need them, please compile the code without that flags.

### Configuration File

It reads the configuration file "ldapPubKeyReader.json" from the following directories:

```
/etc/ldapPubKeyReader.json
/etc/ssh/ldapPubKeyReader.json
./ldapPubKeyReader.json (relative path to binaries location)
```

### Running

This program requires one argument as ldap uid and dumps all of the users public keys to stdout, as SSHD requires. 



