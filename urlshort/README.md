# urlshort

A small web service for handling short urls.

# Next

- Clean up logging in the cmd handler
- Create an admin page behind basic auth to crud url mappings in the boltdb

# Questions

- I can't run `urlshort db add/list` while the db server is running, would it
  be better to open/close the db connection between requests?
  - Seems like no. If we can't connect to the db while another process is, then
    each request would have to wait for the previous open to complete before
    getting to the next.
  - Creates a bottleneck?
