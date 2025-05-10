
on Archlinux
```bash
useradd -mg users -s /usr/bin/zsh test
passwd test

# make it accessible for to the Nginx
chmod 755 /home/test
chmod 755 /home/test/website
```

curl -u youruser:yourpass -T test.txt http://localhost:8800/test.txt


on Debian if login as su from another user 
/usr/sbin/useradd -mg users -s /usr/bin/zsh public

```
GOTOOLCHAIN=auto go mod vendor
GOTOOLCHAIN=auto go run .
```


apt install libpam0g-dev

sudo -u www-data ls /home/$user/website

## Set up the user

chown -R :www-data /home/$user/website
chmod 711 /home/$user