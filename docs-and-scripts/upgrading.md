## Upgrading an Install Done by godojo

godojo's purpose is to get you a working install of DefectDojo. It does not support updating an installation after the initial install is complete. This is by design because of the complexities involved in updating an iron installation, especially one that can be modified post install. Full disclosure: I wrote godojo because maintaining ~2k lines of a Bash script (setup.bash) was getting painful and was remarkably brittle. I prefer and use docker-compose for most of th DefectDojo's I setup and use. However, I wanted an "install on iron" option to remain available to those that perfered that over the other install options.

There's several things godojo does to make updating straightforward:

1. Lays out the installation in a consistent fashion where /opt/dojo/django-DefectDojo (default location) is where the DefectDojo source code is located
2. Setups up a virtualenv for Python3 so that OS Python modules don't interfere with the modules DefectDojo uses
3. Puts the virtualenv inside of /opt/dojo (default location) so all of the app's requirements live in under one directory.

Considerations when updating DefectDojo:

* Backing up the DefectDojo database data  You decide how much redundancy you need. Database backups are not covered here.
  * I typically just do a SQL dump from PostgreSQL or MySQL. Searching for "sql dump" and PostgreSQL or MySQL will give you loads of options.
  * As of this writing, ./manage.py supports "dbbackup" as an option. Details on that are [here](https://github.com/django-dbbackup/django-dbbackup)
* Moving over your DB connection information. This is covered in the example instructions below.
* Moving any customizations from the old version to the new version. Generally these should be none but if you altered the source that godojo installed, you'll need to move those changes to the new version. Good luck with that.

## High-level Upgrade Steps

1. Stop the running instance of DefectDojo services (Celery beat, Celery Worker and Dojo)
2. Relocate the current (to be updated) DefectDojo source
3. Add the new version of DefectDojo source to /opt/dojo/django-DefectDojo (default location)
4. Move over the configurations from the older version of DefectDojo's source
5. Update any Python modules required by the new version
6. Run the commands to update assets and DB migrations to handle any DB changes in the new version
7. Start the DefectDojo services (Celery beat, Celery Worker and Dojo)

## Example Upgrade Walkthrough

This is one way to update an install of DefectDojo. Having been around DefectDojo for 8+years and Linux for 20+, I know there's way more than one way to skin this cat aka do an update.

Note: The example below uses the default paths and options from godojo. If you alter these, you will need to make the corresponding changes to your upgrade process.

(0) Optional but recommended step: Backup the DB data and create a file-level backup of /opt/dojo.

(1) Stop Dojo

```
systemctl stop sysdojo
```

This assumes you use the example scripts provided here. You may have created a service file and would just do `systemctl stop DefectDojo` or could use `kill 15 [pid]`. The key is to make sure that you stop Celery Worker, Celery Beat and DefectDojo (the django app)

(2) Move the current source elsewhere

```
cd /opt/dojo
mv ./django-DefectDojo ./old-dojo
```

Note: It doesn't matter where you move it - you just want the /opt/dojo/django-DefectDojo path to be available.

(3) Add the new version of DefectDojo

In this example, I'll be going from DefectDojo 2.3.1 to 2.4.1 - your version numbers will likely be different and you'll need to update the commands below accordingly.

```
cd /opt/dojo
wget https://github.com/DefectDojo/django-DefectDojo/archive/refs/tags/2.4.1.tar.gz
tar -xzvf 2.4.1.tar.gz
mv django-DefectDojo-2.4.1 django-DefectDojo
rm 2.4.1.tar.gz
ls /opt/dojo
bin  customizations  django-DefectDojo  include  lib  logs  old-dojo  pyvenv.cfg
```
You now have /opt/dojo/django-DefectDojo with the new version of source code.

Note: I used a specific release of DefectDojo - a tarball (.tgz) downloaded from Github above. You could also git clone the DefectDojo repo and checkout the appropriate tag. There's loads of options that will work. The goal is to get the version of the source you want into /opt/dojo/django-DefectDojo

(4) Move over configurations from the older version

When godojo does an install, it puts the provided environment variables into the .env.prod file located at /opt/dojo/django-DefectDojo/dojo/settings. There's a handy symlink to that path at /opt/dojo/customizations.

```
cd /opt/dojo
cp old-dojo/dojo/settings/.env.prod /opt/dojo/django-DefectDojo/dojo/settings/
```

Note: You may have to change the path to the old source code depending on what you did in step (2).

(5) Update the Python modules

Pretty standard stuff. Though the upgrades might take a bit of time.

```
cd /opt/dojo/django-DefectDojo/
/opt/dojo/bin/pip3 install -r requirements.txt --upgrade
```

Note: These updates aren't guarenteed to work without some effort. If new dependencies outside of what pip can do are added to the libraries used by DefectDojo, you may have to do additional installs e.g. `apt install libmysqlclient-dev`.

Also: The reason I use the full path (/opt/dojo/bin/python3) is to use the python installed in the virtualenv without having to activate/deactive it between commands.

(6) Run commands to update Dojo's source

There's a couple of commands to complete the update of the DefectDojo source:

```
cd /opt/dojo/django-DefectDojo/components/
yarn
cd /opt/dojo/django-DefectDojo/
/opt/dojo/bin/python3 ./manage.py collectstatic --noinput
/opt/dojo/bin/python3 ./manage.py migrate
```

(7) Start DefectDojo

```
systemctl start sysdojo
```

You should have 2 Celery, one DefectDojo (python) and a DB running at the least:

```
# pstree
bash-+-postgres---10*[postgres]
     |-pstree
     |-screen---celery-+-3*[celery]
     |                 `-4*[{celery}]
     |-screen---celery
     `-screen---python-+-python---{python}
                       `-11*[{python}]
```

Note: The above command output was generated in a container which explains the lack of other processes in that listing.
