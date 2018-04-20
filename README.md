# Gannet-Market-Api 
[![Build Status](https://travis-ci.com/zricethezav/gannet-market-api.svg?token=jodtRDHhASisqMJ3vY7y&branch=master)](https://travis-ci.com/zricethezav/gannet-market-api)
## Installing and Running the API 
*These instructions assume you have the Go language installed. If you do not, please follow https://golang.org/doc/install*

Building directly from source. Note that you may need to add `$GOPATH/bin` to your `$PATH` in order to
run `gannet-market-api`  
```
go get github.com/zricethezav/gannet-market-api && gannet-market-api
# or 
go get github.com/zricethezav/gannet-market-api && $GOPATH/bin/gannet-market-api
```
or run from docker. `PORT` is up to you:
```
docker run --rm -p PORT:8080 zricethezav/gannet-market-api:latest
```

## Interacting with the API
*The API runs on port 8080*
### Add
The `POST` call adds a produce entry to the database
* **URL**

    /produce

* **Method**
    
    `POST`

* **Body**
    
    `/produce POST` expects a json payload:
    ```
        {"code": <str>, "name": <str>, "price": <float>}
    ```
    * `price` is a positive number up to 2 decimal places
    * `name` is alphanumeric and case insensitive
    * `code` is a code that identifies the produce and case insensitive and
    is sixteen characters long, with dashes separating each four character group
* **Success Response**
    * **Code:** 201 <br />

* **Error Response**

    Error response body is plaintext
    * **Code:** 405 <br />
      **Content:** ` method not allowed`
    * **Code:** 409 <br />
      **Content:** `entry already exists`
    * **Code:** 422 <br />
      **Content:** `malformed request body`

* **Sample Call:**
    ```
    $ curl -X POST -d '{"name":"apple","code":"YRT6-72AS-K736-L4AR", "price": "12.12"}' localhost:8080/produce
    ```
* **Additional Notes**

    Produce codes are unique and if you want to update the price of a produce then you must first delete the produce, then call `/add` with an updated price.
    

### Fetch
The `GET` call retrieves all produce entries in the database
* **URL**

    /produce

* **Method**
    
    `GET`

* **Success Response**
    * **Code:** 200 <br />
    * **Content**: `[...]`
        * example: `[{"Code":"YRT6-72AS-K736-L4AR","Name":"apple","Price":"12.12"},{"Code":"YRT6-72AS-K736-L4AK","Name":"pear","Price":"2.32"}]`

* **Error Response**

    Error response body is plaintext
    * **Code:** 404 <br />
      **Content:** `unable to retrieve entries`
    * **Code:** 405 <br />
      **Content:** `method not allowed`

* **Sample Call:**
    ```
    $  curl -X GET 0.0.0.0:8080/produce
    ```

### Delete 
The `DELETE` call deletes a produce entry from the database based on the url param `code`
* **URL**

    /produce?code=<produce code>

* **Method**
    
    `DELETE`

* **Success Response**
    * **Code:** 204 

* **Error Response**

    Error response body is plaintext
    * **Code:** 404 <br />
      **Content:** `entry does not exist`
    * **Code:** 405 <br />
      **Content:** `method not allowed`
    * **Code:** 422 <br />
      **Content:** `invalid code`

* **Sample Call:**
    ```
    $  curl -X "DELETE" localhost:8080/produce?code=YRT6-72AS-K736-L4ee
    ```

### Deploying
Pushing to master will deploy a build containing the recent changes with the tag `latest`, `master`,
and the Travis build number. Pushing to develop will deploy a building containing develop's changes with the tag
`develop` and the Travis build number.

### Additional Notes:
I enjoyed doing this assignment as I've never set up a CI pipeline from the ground up. This one is simple but I still
learned some useful information about Travis like using the build debugger, how to handle credentials, and I learned the
purpose of `matrix` variables. I didn't actually deploy this to a cloud but if I were to deploy to a cloud provider
I would opt for AWS and make use of their Elastic Beanstalk service.
https://docs.travis-ci.com/user/deployment/elasticbeanstalk/ gives a light walk through on how that process would go.
Test coverage is ~90%. The remaining ~10% untested code is in `func main()` which is responsible for spinning up
the server and defining routes. One final note is that I recognize the popularity of port 8080 and if this were an actual service being deployed somewhere I would change the default port and add the option to configure the port either by passing in an argument or by loading up a yaml config file.

Had to put in some last minute changes after reading https://hackernoon.com/restful-api-designing-guidelines-the-best-practices-60e1d954e7c9 which gives a good rundown on RESTful designs. I think I originally wanted to employ the `/produce`url but got sidetracked by writing tests and reading up on CI. 

One last note on the GET response... I would have liked to test the performance of storing the 'database' entirely as a map rather than a cache + slice. I have a feeling the map route would have been much quicker as it cuts out all code searching to constant time... just lookup the hash.
