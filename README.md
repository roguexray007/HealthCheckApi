# HealthCheckApi
use branch refactoring


                                      ## INTRODUCTION
                                      
a small service to monitor the health of external urls.

The system accept a list of http/https urls and a crawl_timeout(seconds) and frequency(in seconds)
and failure_threshold(count) in JSON format from data.json file

**crawl_timeout** : System will wait for this much time before giving up on the url.    
**frequency** :  System will wait for this much time before retrying again.    
**failure_threshold** :  count of retries possible for that url.

The system iterates over all the urls in the system and try to do a HTTP GET on the URL(wait for the crawl_timeout)
seconds before giving up on the URL. 

###### ROUTES : 

http://localhost:8080/api/healthcheck/addToDB  -> updating db

http://localhost:8080/api/healthcheck/fetchLogs -> fetches records from table logs

###### Cron
time-based job scheduler for running health check after specific intervals (var REFRESHTIME)


