This project is aimed to learn Go.

As a side effect, it will free me from manually download every day the TV 
series I am currently following.

The general design is as follows:
* Get RSS from rlsbb.com
* Compare with list of interesting series
* Find NETLOAD links for 720p downloads
* Get the files in paralel
* Rename the files and copy yhem over to their destination:
  * Series we keep
  * Series we don't keep
* Notify of what is available.