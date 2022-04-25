# JumpCloud Hashing API

## Goal
This API will take in a supplied string to base64 encode and hash using SHA512.

## API Documentation
For more information around API request and responses look at the swagger documentation found [here](https://app.swaggerhub.com/apis/TheTarvis/JumpCloud-Hashing-API/v0.0.1).

## Run
To run the server locally:
```
 go run main.go
```

## Example Usage
Given: A server has just started and no hashes have been saved yet 
First we need to start saving a hash:

```shell
curl --location --request POST 'http://localhost:8080/hash/' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'password=angryMonkey'
```
The expected response should be `1`. Now if we call `/hash/1`, within 5 seconds of the above request, we should get back a `404`.
```shell
curl --location --request GET 'http://localhost:8080/hash/1'
```
If we wait 5 seconds and try the above request again we should get the following response:
```shell
NjQ0MWUxNTgxZWI5ODE0OTczNzU1YzJkMGQwMDJiMTMyYzdlMjk1MmYzYTdmNjkzNjkxNjhmOTQxY2Q4NDQ4MTYzZWFmOGM1NzZhMTFiZDEwZTQxZjMzNTRhMDk5ZDJmMjliNjRmNjY0OTQ5Y2Y0MTVkZWVjYmI2MDNlODFmZWQ=
```

Since we've only saved one password a request to `/stats/` should return:
```json
{
   "total": 1,
   "average": 70
}
```
Note: The average value may be different based on the machines computation time. 

This can be verified by running the following request:
```shell
curl --location --request GET 'http://localhost:8080/stats'
```

And once we're done we can run the shutdown request. If this is done during a hash being saved it will wait til that request is completed.
```shell
curl --location --request GET 'http://localhost:8080/shutdown'
```
### Postman Collection
To run through these test [here](https://www.getpostman.com/collections/65e011d60762744a3f87) is a link to a postman collection with all the requests.

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
### Update `/hash/{id}` response
It returns a 404 for id's that are not in the map. This could be because they're still "processing" or the identifier hasn't been used yet.
### Add automated integration test
Full end-to-end testing right now is manual, I would add some test that verified the full behavior.
### Add authorization to the `/shutdown` call
Right now anyone could shut the server down which isn't great.


## Questions
1. Is there functionality with the returned identifier from `/hash` other than to be able to fetch the stored value?
   1. If there is no other functionality, can I suggest using UUIDs instead to avoid a possible collision with the identifiers.
      1. Answer: Use the base64 encoded string of the SHA512 hash. You can explain in your README.md file why you would recommend using UUIDs instead
2. How long are the encoded strings expected to be available?
   1. Answer: No requirements around expiring encoded strings at this time.
3. Do they need to be resilient on server restart?
   1. Answer: No requirement to keep strings resilient at this time.
4. To allow for backwards compatibility, is it okay to version in the API endpoints? For example, `/v1/hash`?
   1. Answer: No requirement to allow for backwards compatibility. You can explain in your README.md file if you use /v1/
5. I see a requirement to not use packages outside https://pkg.go.dev/std#stdlib, is this a hard requirement or can I use gin, https://github.com/gin-gonic/gin, for the webserver?
   1. Answer: It's a hard requirement
6. Any requirements around payload validation? For example: password field not supplied in form data.
   1. Answer: handle it how you would if it were going in production
7. For the `/stats` response body, is the total for both `/hash/` and `/hash/42` request? Or just `/hash`
   1. Answer: only for /hash