# /usr/share/webdav/templates
all:
	go build  -o webdav ./src
	mkdir -p /srv/bin
	mv webdav /srv/bin/
	cp webdav.service /etc/systemd/system/
	systemctl enable webdav
	systemctl start webdav