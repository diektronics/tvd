This project is aimed to learn Go.

As a side effect, it will free me from manually download every day the TV 
series I am currently following.

The general design is as follows:
* Get RSS from rlsbb.com
* If if has changed, compare with list of interesting series and check if
  new eps are available
* Find NETLOAD links for 720p downloads
* Get the files in paralel using a pool of workers
  * Rename the files and copy them over to their destination:
    * Series we keep
    * Series we don't keep
    * Notify of what is available.
* Sleep 20'

Interesting series will be kept in a DB:
  Series name, last episode downloaded, location
  Check https://github.com/Go-SQL-Driver/MySQL/