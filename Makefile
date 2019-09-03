VERSION= 1.0.0
LDFLAGS= -ldflags "-X main.version=${VERSION}"
RPM_BUILD_PATH= ~/rpmbuild
RPM_BUILD_ROOT= ${RPM_BUILD_PATH}/BUILDROOT

default: build

run: build
	cd monitor; ./rsa-nw-monitor

depends:
	go get -d ./...

build: depends
	cd monitor; go build $(LDFLAGS) -o rsa-nw-monitor

clean:
	rm -f monitor/rsa-nw-monitor
	rm -Rf ${RPM_BUILD_PATH}

install: build
	mkdir -p ${RPM_BUILD_ROOT}
	mkdir -p ${RPM_BUILD_ROOT}/usr/local/bin/
	mkdir -p ${RPM_BUILD_ROOT}/etc/rsa-nw-monitor/
	mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system/
	cp monitor/rsa-nw-monitor ${RPM_BUILD_ROOT}/usr/local/bin
	cp conf/rsa-nw-monitor.conf ${RPM_BUILD_ROOT}/etc/rsa-nw-monitor
	cp conf/rsa-nw-monitor.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system

rpm: install
	mkdir -p ${RPM_BUILD_PATH}/SPECS ${RPM_BUILD_PATH}/RPMS ${RPM_BUILD_PATH}/SOURCES
	cp conf/rpmbuild/SPECS/rsa-nw-monitor.spec ${RPM_BUILD_PATH}/SPECS
	sed -i 's/%VERSION%/${VERSION}/' ${RPM_BUILD_PATH}/SPECS/rsa-nw-monitor.spec
	cp LICENSE ${RPM_BUILD_PATH}/SOURCES/license
	rpmbuild -ba ${RPM_BUILD_PATH}/SPECS/rsa-nw-monitor.spec
