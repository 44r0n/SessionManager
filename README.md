# Sessionmanager

A little MySQL database to manage users.

[![Build status](https://travis-ci.org/44r0n/Sessionmanager.svg)](https://travis-ci.org/44r0n/Sessionmanager)
[![GitHub release](https://img.shields.io/github/release/44r0n/sessionmanager.svg)](https://github.com/44r0n/Sessionmanager/releases)
[![Libraries.io for GitHub](https://img.shields.io/librariesio/github/44r0n/Sessionmanager.svg)]()

## Getting Started

### Prerequisites

You need [Docker](https://www.docker.com/) or [MySQL](https://www.mysql.com/) installed and check that [mytap](https://github.com/theory/mytap) is a submodule and it is in the project folder.

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

This file just sets up the database. If you changed the configuration previusly you must change it in this file.

## Running the tests

To run the tests just execute de `initsql` script with the `-t` flag:
~~~
~/path_to_the_project$ ./initsql.sh -t
~~~

## Usage

There are some main functions or stored procedures to make it easy to use.

~~~
CreateUser ('username','password','mail@mail.com');
~~~                              

Creates a user, with the username and hashes the password.

~~~
Login ('username', 'password');
~~~

Checks if the useer exists with the given username and password. If so, returns an identificative token, if not returns a 0.

~~~
TokenCleaner ();
~~~

Cleans the tokens within an predefined time, default time is 30 mins. You can configure this time by inserting it in the table `configuration` setting the `name` to `Expiration time` and `value` to the time in minutes you wish.

~~~
Logout ('token');
~~~

Deletes the given token  in the `user_tokens` so the user is no more logged in the sistem.


## Built With

*   [MySQL](https://www.mysql.com/) - The databse engine.
*   [mytap](https://github.com/theory/mytap) - The test framework.

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
