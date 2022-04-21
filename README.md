# JumpCloud Hashing API

## Goal
This API will take in a supplied string to base64 encode and hash using SHA512.

## API Documentation
swagger?
```json
{
  
}
```

## Installation

## Run

## Test
To run unit test for the application execute the following command:
```
 go test ./... -cover
```
## Known Issues
1. Currently, the hashes are being stored in local memory and given enough uptime and request we will run out of memory.
   1. Solution here would be to either limit the time which the data lives, save the data to a database with some caching to help minimize fetch times.
2. On restart all hashes and their corresponding identifier will be lost.
   1. Persistent storage could be used to fix this.
3. 

## Enhancements

### Using GIN for webserver 
I've found this library to be one of the best ways to develop webservers using Go. More information can be found [here](https://github.com/gin-gonic/gin).
### Using versioning in REST endpoints
To support backwards compatability I recommend adding some form of versioning to the endpoints. For example adding a `/v1` to the endpoints so in the future if there is a breaking change we can create a new endpoint under `/v2`.
### Fix trailing slash routing issue
Right now we have to register `/hash` and `/hash/` otherwise a request to `/hash` is routed as `GET` request, and I'm not 100% sure why

## Questions
1. ~~Is there functionality with the returned identifier from `/hash` other than to be able to fetch the stored value?~~
   1. ~~If there is no other functionality, can I suggest using UUIDs instead to avoid a possible collision with the identifiers.~~
2. ~~How long are the encoded strings expected to be available?~~ 
3. ~~Do they need to be resilient on server restart?~~
4. ~~To allow for backwards compatibility, is it okay to version in the API endpoints? For example, `/v1/hash`?~~
5. ~~I see a requirement to not use packages outside https://pkg.go.dev/std#stdlib, is this a hard requirement or can I use gin, https://github.com/gin-gonic/gin, for the webserver?~~
6. Any requirements around payload validation? For example: password field not supplied in form data. 