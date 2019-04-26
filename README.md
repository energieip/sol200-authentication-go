Authentication: Service Management for configuration and data
=============================================================

Authentication is responsible for:
* managing users
* creating tokens

Build Requirement: 
* golang-go > 1.9
* glide
* devscripts
* make

Run dependancies:
* rethindkb

To compile it:
* GOPATH needs to be configured, for example:
```
    export GOPATH=$HOME/go
```

* Install go dependancies:
```
    make prepare
```

* To clean build tree:
```
    make clean
```

* Multi-target build:
```
    make all
```

* To build x86 target:
```
    make bin/authentication-amd64
```

* To build armhf target:
```
    make bin/authentication-armhf
```
* To create debian archive for x86:
```
    make deb-amd64
```

* To install debian archive on the target:
```
    scp build/*.deb <login>@<ip>:~/
    ssh <login>@<ip>
    sudo dpkg -i *.deb
```

* To add user on target:
```
   /usr/local/bin/add-user  -c /etc/energieip-sol200-authentication/config.json -u admin -p admin
```

For development:
* recommanded logger: *rlog*
* For database/network connection: use *common-components-go* library
