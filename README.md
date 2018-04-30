# Sample API
[![Build Status](https://travis-ci.com/zricethezav/sample-api.svg?token=jodtRDHhASisqMJ3vY7y&branch=master)](https://travis-ci.com/zricethezav/sample-api)
## Installing and Running the API 
*These instructions assume you have the Go language installed. If you do not, please follow https://golang.org/doc/install*

Building directly from source. Note that you may need to add `$GOPATH/bin` to your `$PATH` in order to
run `sample-api`  
```
go get github.com/zricethezav/sample-api && sample-api
# or 
go get github.com/zricethezav/sample-api && $GOPATH/bin/sample-api
```
or run from docker. `PORT` is up to you:
```
docker run --rm -p PORT:8080 zricethezav/sample-api:latest
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
