###############################################################################
# Spec file for rsa-nw-monitor
################################################################################
# Configured to be built by non-root user
################################################################################
#
# Build with the following syntax:
# rpmbuild -bb rsa-nw-monitor.spec
#
Summary: Monitor Service providing a REST API interface to retrieve RSA Netwitness EndPoint Host Information
Name: rsa-nw-monitor
Version: %VERSION%
Release: 1
License: Apache
Group: Utilities
Packager: Helmut Wahrmann
BuildRoot: ~/rpmbuild/

%description
The Monitor Service can be configured via Command line or in a config file.
The config file can be specified as part of the command line.
If not present the default config in /etc/rsa-nw-monitor/rsa-nw-monitor.conf
will be used.

Netwitness Monitor Service provides a REST API to retrieve the Risk Score of a JoinHostPort
for monitoring purposes.

%prep
echo "BUILDROOT = $RPM_BUILD_ROOT"
mkdir -p $RPM_BUILD_ROOT
mv $RPM_BUILD_ROOT/../usr/ $RPM_BUILD_ROOT/
mv $RPM_BUILD_ROOT/../etc/ $RPM_BUILD_ROOT/

%files
%attr(0744, root, root) /usr/local/bin/*
%config %attr(0644, root, root) /etc/rsa-nw-monitor/*
%attr(0644, root, root) /usr/lib/systemd/system/rsa-nw-monitor.service

%post
################################################################################
# Set up a sybobilc link to our new service                                    #
################################################################################
cd /etc/systemd/system/multi-user.target.wants

if [ ! -e rsa-nw-monitor.service ]
then
   ln -s /usr/lib/systemd/system/rsa-nw-monitor.service
fi

%postun
# remove installed files and links
rm /etc/systemd/system/multi-user.target.wants/rsa-nw-monitor.service

%clean
rm -rf $RPM_BUILD_ROOT/usr
rm -rf $RPM_BUILD_ROOT/etc

%changelog
* Tue Sep 03 2019 Helmut Wahrmann <helmut.wahrmann@rsa.com>
  - Inital release
