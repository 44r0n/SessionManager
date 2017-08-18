# Sessionmanager

A little MySQL database to manage users.

[![Build status](https://travis-ci.org/44r0n/SessionManager.svg)](https://travis-ci.org/44r0n/SessionManager)
[![GitHub release](https://img.shields.io/github/release/44r0n/sessionmanager.svg)](https://github.com/44r0n/Sessionmanager/releases)
[![Libraries.io for GitHub](https://img.shields.io/librariesio/github/44r0n/Sessionmanager.svg)]()

## Getting Started

### Prerequisites

You need [Docker](https://www.docker.com/) or [MySQL](https://www.mysql.com/) installed and check that [mytap](https://github.com/theory/mytap) is a submodule and it is in the project folder.
Also you need [Go](https://golang.org) installed with the following dependencies:
-   [httprouter](github.com/julienschmidt/httprouter)
-   [MySQL driver](github.com/go-sql-driver/mysql)
-   [Configuration](github.com/tkanos/gonfig)
-   [JWT](github.com/dgrijalva/jwt-go)
-   [BCrypt](golang.org/x/crypto/bcrypt)
-   [GoConvey](github.com/smartystreets/goconvey/convey)

### Installing

If you are using [Docker](https://www.docker.com/) go to the project folder and execute:

~~~
~/path_to_the_project$ ./initDocker.sh
~~~

In this file you can find the configuration of the [MySQL](https://www.mysql.com/) server, feel free to modify it.
Then execute:
~~~
~/path_to_the_project$ ./initsql.sh
~~~

This file just sets up the database by default in the [Docker](https://www.docker.com/) machine. If you changed the configuration previously you must change it in this file or specify in the second parameter:
~~~
~/path_to_the_project$ ./initsql.sh '' mysql.ip
~~~

Installing dependencies:
~~~
~/path_to_the_project$ ./installDependencies.sh
~~~

Once the database is set up and the dependecies installed, build and run the `server.go` file.
~~~
~/path_to_the_project$ go build server.go
~/path_to_the_project$ ./server
~~~

## Running the tests

To run the tests just execute de `initsql` script with the `-t` flag:
~~~
~/path_to_the_project$ ./initsql.sh -t
~~~

Like the previous step this tries to execute the tests in the [Docker](https://www.docker.com/) machine. You can change the ip by parameter:
~~~
~/path_to_the_project$ ./initsql.sh -t mysql.ip
~~~

To test the server integrated with the database run the tests with `-database`flag:
~~~
~/path_to_the_project$ go test ./... -v -database
~~~

## API Usage

##### Register User

Registers a json given user.

*   **URL**

    /Register

*   **Method:**

    `POST`

*   **URL Params**

    None

*   **Data Params**

    ~~~json
    {
    	"username":"Bob Smith",
    	"email":"mail@mail.com",
    	"password":"secretPassword"
    }
    ~~~

*   **Success Response:**

*   **Code:** 201 CREATED<br />
      **Content:** No content

*   **Error Response:**

*   **Code:** 409 CONFLICT <br />
    **Content:**
    ~~~json
    {
    	"error": {
    		"status": 409,
    		"error": "FIELDS_REPEATED",
    		"description": "One or more fields already exist",
    		"fields": {
    			"email": "An account already exists with this email",
    			"username": "An account already exists with this username"
    		}
    	}
    }
~~~

* **Sample Call:**

  ~~~javascript
    $.ajax({
      url: "/Register",
      dataType: "json",
      type : "POST",
      201 : function(r) {
        console.log(r);
      }
    });
  ~~~
----
##### Login User

  Logs in with a json given user

*   **URL**

      /Login

*   **Method:**

      `POST`

*   **URL Params**

      None

*   **Data Params**

      ~~~json
      {
    	"username":"Bob Smith",
    	"password":"secretPassword"
      }
      ~~~

*   **Success Response:**

    *   **Code:** 200 OK<br />
        **Content:**
        ~~~json
        {
        	"response": {
        		"status": "OK",
        		"token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImQxZmY1MDE4LTgyYjgtMTFlNy1hMGQ1LTAyNDJhYzExMDAwMiJ9.vwm1JhqfPhI6Vj-A4BubhaGyrSpGSKYQ246JSPQqd6lZBrkpuDrpYem79baUJPQDXoduw14j6x26KkK8wprzTg",
        		"error": ""
        	}
        }
        ~~~

*   **Error Response:**

*   **Code:** 404 NOT FOUND <br />
    **Content:**
      ~~~json
      {
      	"response": {
      		"status": "Incorrect user or password",
      		"token": "",
      		"error": ""
      	}
      }
      ~~~

*   **Sample Call:**
    ~~~javascript
      $.ajax({
        url: "/Register",
        dataType: "json",
        type : "POST",
        200 : function(r) {
          console.log(r);
        }
      });
    ~~~

----
##### Validate Token

Validates a given token.

*   **URL**

      /Token/isValid

*   **Method:**

      `POST`

*   **URL Params**

      None

*   **Data Params**

      None

*   **Headers**

      `"Authorization":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImQxZmY1MDE4LTgyYjgtMTFlNy1hMGQ1LTAyNDJhYzExMDAwMiJ9.vwm1JhqfPhI6Vj-A4BubhaGyrSpGSKYQ246JSPQqd6lZBrkpuDrpYem79baUJPQDXoduw14j6x26KkK8wprzTg"`

*   **Success Response:**

    *   **Code:** 204 NO CONTENT<br />
        **Content:**
          None

*   **Error Response:**

*   **Code:** 404 NOT FOUND <br />
    **Content:**
      None

*   **Sample Call:**
    ~~~javascript
      $.ajax({
        url: "/Register",
        dataType: "json",
        type : "POST",
        204 : function(r) {
          console.log(r);
        }
      });
    ~~~    
----
##### Logout User

  Logs out with a json given user

*   **URL**

      /Logout

*   **Method:**

      `POST`

*   **URL Params**

      None

*   **Data Params**

      ~~~json
      {
    	"username":"Bob Smith",
    	"password":"secretPassword"
      }
      ~~~

*   **Success Response:**

    *   **Code:** 204 NO CONTENT<br />
        **Content:**
          None

*   **Error Response:**

*   **Code:** 404 NOT FOUND <br />
    **Content:**
      ~~~json
      {
      	"response": {
      		"status": "Incorrect user or password",
      		"token": "",
      		"error": ""
      	}
      }
      ~~~

*   **Sample Call:**
    ~~~javascript
      $.ajax({
        url: "/Register",
        dataType: "json",
        type : "POST",
        201 : function(r) {
          console.log(r);
        }
      });
    ~~~

## Built With

*   [MySQL](https://www.mysql.com/) - The database engine.
*   [mytap](https://github.com/theory/mytap) - The test framework.
-   [httprouter](github.com/julienschmidt/httprouter) - Routing.
-   [MySQL driver](github.com/go-sql-driver/mysql) - MySqlDriver.
-   [Configuration](github.com/tkanos/gonfig) - Configuration stuff.
-   [JWT](github.com/dgrijalva/jwt-go) - Building JWT.
-   [BCrypt](golang.org/x/crypto/bcrypt) - BCrypt security.
-   [GoConvey](github.com/smartystreets/goconvey/convey) - Makes testing quick & easy.

## Contributing

Please read [CONTRIBUTING.md](https://github.com/44r0n/Sessionmanager/blob/master/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/44r0n/Sessionmanager/tags).

## Integrations
-   [ ] Google
-   [ ] Github
-   [ ] Twitter
-   [ ] Facebook

## Authors

*   **Aarón Sánchez Navarro** - *Initial work* - [Sessionmanager](https://github.com/44r0n/Sessionmanager)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
