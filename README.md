# product_recommendation

## Go version

`1.24`

## How to launch the service in localhost
1. build docker image
```
make docker-build
```
2. (Optional) Customize the `config.yml` file if needed for your environment.
3. Start the service using docker-compose:
```
docker-compose up -d
```

## Postman Collection

To get started quickly with the API, simply import the provided `product_recommendation.postman_collection.json` file from this repository into Postman.

## Go Test Usage

### Overview

This repository contains a Go project that can be tested using the `make test` command.

### Usage

To run the tests, simply execute the following command:
```
make test
```

## What does this repository do?

This repository provides a backend service built with Go (version 1.22 or above), MySQL, and Redis. It implements a user authentication and product recommendation system, suitable for backend technical interviews. The service includes the following features:

1. **User Registration**
   - Users register with an email address and password.
   - Passwords must be 6-16 characters long, contain at least one uppercase letter, one lowercase letter, and one special character from ()[]{}<>+-*/?,.:;"'_\|~`!@#$%^&=.

2. **Email Verification**
   - Supports either (or both) of the following methods:
     - Email verification via a link.
     - Email verification via a code sent to the user's email, which the user submits to verify.
   - The email sending function is a stub and does not actually send emails.

3. **User Login**
   - Users can log in using their email and password.

4. **GET /recommendation API**
   - Authenticated users can access a product recommendation list.
   - The recommendation data is cached in Redis with a 10-minute expiration to handle high load (expected 300 requests per minute).
   - The underlying database query for recommendations is intentionally slow (3 seconds) to demonstrate caching effectiveness.

This project does not include frontend implementation. Please refer to the rest of this README for setup and usage instructions.
