# Instagram Clone API

---

## An Instagram clone app with Golang using only standard packages and mongodb-driver package.

## Installation

1. Clone this repo
2. Cd to repo directory and create a secrets.go file [I am giving out secrets for easy testing]
```
package main
var secretDbURI string = "mongodb+srv://samplemongouri.net/"
var secretHashKey string = "a secret string"
var secretPort string = ":9000" //port to listen

```

3. build and run `./instagram-clone-appointy`
4. alternative - use docker to build and run.

## Approach

Made a tree based routing. used ces encryption and decryption for password hashing. mongodb limit and skip for pagination.

## Routes

Standard routes as given in the task doc. All routes work as expected.

for seeking all posts pagination done by simply taking out limit and skip params from url `/users/posts/:id?limit=10&skip=6`
