# Gannet-Market-Api 
[![Build Status](https://travis-ci.com/zricethezav/gannet-market-api.svg?token=jodtRDHhASisqMJ3vY7y&branch=master)](https://travis-ci.com/zricethezav/gannet-market-api)
## Running 
Building directly from source. Note that you may need to add `$GOPATH/bin` to your `$PATH` in order to
run `gannet-market-api`  
```
go get -u github.com/zricethezav/gannet-market-api && gannet-market-api
# or 
go get -u github.com/zricethezav/gannet-market-api && $GOPATH/bin/gannet-market-api
```
or run from docker. `PORT` is up to you):
```
docker run --rm -p PORT:8080 zricethezav/gannet-market-api:latest
```

## Interacting with the API
### Add
The `/add` call adds a produce entry to the database
* **URL**

    /add

* **Method**
    
    `POST`

* **Body**
    
    `/add` expects a json payload:
    ```
        {"code": <str>, "name": <str>, "price": <float>}
    ```
* **Success Response**
    * **Code:** 201 <br />
      **Content:** `{ id : 12 }`

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
    $ curl -X POST -d '{"name":"apple","code":"YRT6-72AS-K736-L4AR", "price": "12.12"}' localhost:8080/add
    ```
    

### Fetch
The `/fetch` call retrieves all produce entries in the database
* **URL**

    /fetch

* **Method**
    
    `GET`

* **Success Response**
    * **Code:** 200 <br />
      **Content:** `[...]`

* **Error Response**

    Error response body is plaintext
    * **Code:** 404 <br />
      **Content:** `unable to retreive entries`
    * **Code:** 405 <br />
      **Content:** `method not allowed`

* **Sample Call:**
    ```
    $  curl -X GET 0.0.0.0:8080/fetch
    ```

### Delete 
The `/delete` call deletes a produce entry from the database based on the url param `code` 
* **URL**

    /delete

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
    $  curl -X "DELETE" localhost:8080/delete?code=YRT6-72AS-K736-L4ee
    ```

