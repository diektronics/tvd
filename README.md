This project is aimed to learn Go.

As a side effect, it will free me from manually download every day the TV 
series I am currently following.

The general design is as follows:
* Get RSS from rlsbb.com
* If it has changed, compare with list of interesting series and check if
  new episodes are available
* Find NETLOAD links for 720p downloads
* Get the files in paralel using a pool of workers
  * Rename the files and copy them over to their destination:
    * Series we keep
    * Series we don't keep
    * Notify of what is available.
* Sleep 20'

Interesting series will be kept in a DB with a table containing: Series name, last episode downloaded, location.

Commands to generate DB:
      $ mysql -u root -p 
      #> grant all on tvd.* to tvd identified by 'tvd';
      #> flush privileges;
      #> quit;
      $ mysql -u tvd -p
      #> create database tvd;
      #> use tvd;
      #> create table series (serie_id int NOT NULL AUTO_INCREMENT, name varchar(255) NOT NULL, latest_ep varchar(10) DEFAULT NULL, location varchar(255) NOT NULL, PRIMARY KEY (serie_id));
      #> create unique index unique_series_name on series (name);

TODO:
* Do not update a new ep until it has actually been downloaded.
* Allow for recovering currently downloading eps and eps in queue.