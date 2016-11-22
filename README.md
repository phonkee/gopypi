# gopypi

Gopypi is private pypi repository implemented in go language. Gopypi can store your private packages and support
`pip install`.
Gopypi supports distutils/setuptools setup.py uploading of packages.
And it does more thant this.
Gopypi supports management of active users, maintainers of packages and lot more.

**This project is under development and should not be used in production for now. I am working hard to 
implement all features and bring alpha version in short time.**

## Features

If you install gopypi you will get out of the box following features:

* auth module - maintain users with access permissions, maintainers access to packages
* packages - uploading of packages, versions, files (directly from setup.py)
* Licenses - list of licenses for uploaded packages
* Clean admin interface - Simple but powerful modern SPA admin interface (written in vue.js)
* Embedded static data - all static files are embedded directly into binary
 
Gopypi has some optional features that can be enabled in administration, such as:

* download stats - stats about downloaded packages
* automatically assign maintainers - reads information from python packages and assigns users

## Admin

Gopypi has modern SPA admin interface where you can browse information about packages, maintain user credentials etc.
Currently only admin can login into this administration, but that will change in near future, so all active users
can log into admin and see information about packages they maintain.

## Command line interface
gopypi is single binary which has multiple subcommands, you can see them by typing `gopypi -h`

## Installation

The easiest way is to download precompiled binaries from github repository releases.
Current version ccan be found here:

https://github.com/phonkee/gopypi/releases/tag/0.5.3

Otherwise you can type

    go get github.com/phonkee/gopypi/...

After you have installed gopypi you need to create config file. You can use gopypi command `makeconfig` that helps
you with generating new configuration interactively.

    `./gopypi makeconfig`

### Database

Gopypi currently supports postgres, mysql databases. To correctly setup database dsn refer to
 
* postgres - https://godoc.org/github.com/lib/pq
* mysql - https://github.com/go-sql-driver/mysql#dsn-data-source-name

Now it's time to run database migrations:

    `./gopypi migrate --config gopypi.conf`
    
This will apply all database migrations automaticaly. For every new gopypi release migrate command is needed to be run.

After you have applied all database migrations it's time to create new admin user, so you can access admin interface:

    `./gopypi createadmin --config gopypi.conf`

After you have succesfully created admin user you can now run server.

    `./gopypi runserver --config gopypi.conf`


## Future features

Gopypi has following features planned:

* list of package classifiers with stats
* only admin can login into admin interface
* enable registration from `python setup.py register`
* email notifications about package changes
* classifiers - create page in admin groupping packages by classifier

## Upload package

For instructions how to write `setup.py` you can find more information in admin interface `HOWTO`. 

## Screenshots

Here are some screenshots from gopypi admin

#### Admin dashboard
![Dashboard](/docs/dashboard.png?raw=true "Dashboard")

#### Admin user list
![User list](/docs/users.png?raw=true "User list")

#### Admin add user
![Add user](/docs/adduser.png?raw=true "Add user")

#### Admin package detail
![Package detail](/docs/packagedetail.png?raw=true "Package detail")

#### Admin download stats
![Download stats](/docs/dstats.png?raw=true "Download stats")

# Contributors
You are welcome!

# Author
phonkee