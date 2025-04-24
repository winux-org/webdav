

```bash
useradd -mg users -s /usr/bin/zsh test
passwd test

# make it accessible for to the Nginx
chmod 755 /home/test
chmod 755 /home/test/website
```

curl -u youruser:yourpass -T test.txt http://localhost:8800/test.txt


```
GOTOOLCHAIN=auto go mod vendor
GOTOOLCHAIN=auto go run .
```


apt install libpam0g-dev
