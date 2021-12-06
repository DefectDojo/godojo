## Bonus Docs and Scripts

godojo's primary job is to get a working install of DefectDojo onto a VM. How you want to run that install is a matter of opinion.  Are you running for yourself? Just experimenting/playing around with DefectDojo? Planning on having a team of people use DefectDojo? Is the install on a cloud provider, a corporate LAN or just on your laptop?

So, there's not **one way** to run DefectDojo after godjo is done. It realy depends on how you want to use DefectDojo and where you install it.

Just to get you started, I've added these two simple scripts to this repo. They will start DefectDojo (the application and the Celery services) for you. They assume the database you chose (MySQL or PostgreSQL) is already running.

Consider this a good starting place to figuring out the install that works for your situation.

**Upgrading instructions** for DefectDojo installed by godojo is located [here](https://github.com/DefectDojo/godojo/tree/master/docs-and-scripts/upgrade.md)

### Using the scripts

To use these as is, copy them to your install, make them executable and run them like:

```
# cd /opt/dojo
# cp /path/where/you/put/the/scripts/dojo-start ./
# cp /path/where/you/put/the/scripts/dojo-stop ./
# chmod 775 dojo-start dojo-stop
# ./dojo-start

  [ after a while...]

# ./dojo-stop
```

Many people like to put Nginx or Apache in front of DefectDojo to handle TLS and other standard web server configurations. These scripts allow you to do that by simply using a reverse proxy of http://127.0.0.1:8000.  You can read more about setting up a reverse proxy at:

* https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/
* https://httpd.apache.org/docs/2.4/howto/reverse_proxy.html

As always, additional ways to run DefectDojo are gladly accepted as PRs to the Community-Contribs repo at https://github.com/DefectDojo/Community-Contribs, PRs to the godojo/scripts directory or published to your blog, etc.

Enjoy!
