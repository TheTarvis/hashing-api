# JumpCloud Hashing API

## Goal
This API will take in a supplied string to base64 encode and hash using SHA512.

## API Documentation
For more information around API request and responses look at `swagger.yml`

## Run
To run the server locally:
```
 go run main.go
```
## Test
To run unit test for the application execute the following command:
```
 go test ./... -cover
```

There is also a file `loasttest/loadtest.go` that you can execute to run a load test against the local webserver. This needs some TLC but it works...
## Known Issues
1. Currently, the hashes are being stored in local memory and given enough uptime and request we will run out of memory.
   1. Solution here would be to either limit the time which the data lives, save the data to a database with some caching to help minimize fetch times.
2. On restart all hashes and their corresponding identifier will be lost.
   1. Persistent storage could be used to fix this.
3. Identifiers will run out of unique values since they're auto-incrementing integers. 

## Enhancements
### Break the state file up
There is too much going on in here and some separation of responsibilities and its data would be a good idea.
### Using GIN for webserver 
I've found this library to be one of the best ways to develop webservers using Go. More information can be found [here](https://github.com/gin-gonic/gin).
### Using versioning in REST endpoints
To support backwards compatability I recommend adding some form of versioning to the endpoints. For example adding a `/v1` to the endpoints so in the future if there is a breaking change we can create a new endpoint under `/v2`.
### Fix trailing slash routing issue
Right now we have to register `/hash` and `/hash/` otherwise a request to `/hash` is routed as `GET` request, and I'm not 100% sure why
### Use prebuilt libraries for metrics
The `/stats` endpoint is pretty basic, I would suggest adding [prometheus](https://prometheus.io/docs/guides/go-application/) or any other library.
### Add logging
There is minimal logging here, I would suggest we add more with more details like request bodies.
### Remove hardcoded values to env vars
For example, port is hardcoded to `:8080` but this should be set at runtime.
### Better error handling
Right now most errors are logged but improvements could be made around this.

## Questions
1. ~~Is there functionality with the returned identifier from `/hash` other than to be able to fetch the stored value?~~
   1. ~~If there is no other functionality, can I suggest using UUIDs instead to avoid a possible collision with the identifiers.~~
2. ~~How long are the encoded strings expected to be available?~~ 
3. ~~Do they need to be resilient on server restart?~~
4. ~~To allow for backwards compatibility, is it okay to version in the API endpoints? For example, `/v1/hash`?~~
5. ~~I see a requirement to not use packages outside https://pkg.go.dev/std#stdlib, is this a hard requirement or can I use gin, https://github.com/gin-gonic/gin, for the webserver?~~
6. Any requirements around payload validation? For example: password field not supplied in form data.
7. For the `/stats` response body, is the total for both `/hash/` and `/hash/42` request? Or just `/hash`