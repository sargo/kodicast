# some helpful shortcuts

build:
	go install github.com/sargo/kodicast

fmt:
	go fmt . ./apps ./apps/youtube ./apps/youtube/mp ./config ./log ./server

run: build
	../../bin/kodicast

install:
	cp ../../bin/kodicast /usr/local/bin/kodicast.new
	mv /usr/local/bin/kodicast.new /usr/local/bin/kodicast
	if ! egrep -q "^kodicast:" /etc/passwd; then useradd -s /bin/false -r -M kodicast -g audio; fi
	mkdir -p /var/local/kodicast
	chown kodicast:audio /var/local/kodicast
	cp $(CURDIR)/kodicast.service /etc/systemd/system/kodicast.service
	systemctl enable kodicast

remove:
	rm -f /usr/local/bin/kodicast
	if egrep -q "^kodicast:" /etc/passwd; then userdel kodicast; fi
	# rm -rf /var/local/kodicast # this removes configuration files
	systemctl disable kodicast
	rm -f /etc/systemd/system/kodicast.service
