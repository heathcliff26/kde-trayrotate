%global debug_package %{nil}

Name:           kde-trayrotate
Version:        0
Release:        %autorelease
Summary:        Quicksetting for auto rotation in the system tray

%global package_id io.github.heathcliff26.%{name}

License:        Apache-2.0
URL:            https://github.com/heathcliff26/kde-trayrotate
Source:         %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires: golang >= 1.27
Requires: libkscreen

%global _description %{expand:
Quickly toggle the KDE auto rotate setting for your screen from the system tray.}

%description %{_description}

%prep
%autosetup -n kde-trayrotate-%{version} -p1

%build
export RELEASE_VERSION="%{version}-%{release}"
make build

%install
install -D -m 0755 bin/%{name} %{buildroot}/%{_bindir}/%{name}
install -D -m 0644 packages/%{package_id}.desktop %{buildroot}/%{_datadir}/applications/%{package_id}.desktop
install -D -m 0644 packages/%{package_id}.png %{buildroot}/%{_datadir}/icons/hicolor/512x512/apps/%{package_id}.png
install -D -m 0644 packages/%{package_id}.svg %{buildroot}/%{_datadir}/icons/hicolor/scalable/apps/%{package_id}.svg
install -D -m 0644 %{package_id}.metainfo.xml %{buildroot}/%{_datadir}/metainfo/%{package_id}.metainfo.xml

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}
%{_datadir}/applications/%{package_id}.desktop
%{_datadir}/icons/hicolor/512x512/apps/%{package_id}.png
%{_datadir}/icons/hicolor/scalable/apps/%{package_id}.svg
%{_datadir}/metainfo/%{package_id}.metainfo.xml

%changelog
%autochangelog
