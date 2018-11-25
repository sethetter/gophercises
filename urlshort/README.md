# urlshort

A small web service for handling short urls.

# Next

- Accept a yaml file arg and use that mapping instead of the default mapping
- Create a JSONHandler func and optionally accept a json file arg (allow only one)
- Create a handler that reads from boltdb, behind a flag (db name? file location? idk)
- Create an admin page behind basic auth to crud url mappings in the boltdb
