# fetchercise
An exercise for Fetch ! 

# Prereqs 

You will need go, git, and something to request 
If you use postman you can use the following collection to test all the requests that were mentioned in the pdf

https://www.getpostman.com/collections/48d8e08a3dc3bd356ba8

If you don't have go here's a link to download
https://go.dev/doc/install

# Caveats

I think I immediately fail because I don't match the results but I didnt understand why the balance isnt updated as the posts are updated and by the time I realized that 
I felt like I've spent too much time on this (about a full working day) too much from a standpoint of if I take too long to do this than maybe I don't fit in. Thank you for giving me this challenge

# Run 

from any folder 
I'm not sure if you need to run go mod init it should be included in the pull but just incase 

	go mod init fetchercise/interview

	git clone git@github.com:edheadonfire/fetchercise.git
	cd fetchercise
	go get . 
	go run . 
	

	

# Curls

## POSTS

	curl --location --request POST 'localhost:8080/pay' \
	--header 'Content-Type: application/json' \
	--data-raw '{
    	"payer":"DANNON",
    	"points":1000,
    	"timestamp":"2020-11-02T14:00:00Z"
	}'
	
	curl --location --request POST 'localhost:8080/pay' \
	--header 'Content-Type: application/json' \
	--data-raw '{ "payer": "UNILEVER", "points": 200, "timestamp": "2020-10-31T11:00:00Z" }'

	curl --location --request POST 'http://localhost:8080/pay' \
	--header 'Content-Type: application/json' \
	--data-raw '{
	    "payer": "DANNON",
	    "points": -200,
	    "timestamp": "2020-10-31T15:00:00Z"
	}'

	curl --location --request POST 'http://localhost:8080/pay' \
	--header 'Content-Type: application/json' \
	--data-raw '{
	    "payer": "MILLER COORS",
	    "points": 10000,
	    "timestamp": "2020-11-01T14:00:00Z"
	}'

	curl --location --request POST 'http://localhost:8080/pay' \
	--header 'Content-Type: application/json' \
	--data-raw '{ "payer": "DANNON", "points": 300, "timestamp": "2020-10-31T10:00:00Z" }
	'

## Get

	curl --location --request GET 'http://localhost:8080/checkBalances'

## Put

	curl --location --request PUT 'http://localhost:8080/spend' \
	--header 'Content-Type: application/json' \
	--data-raw '{
	    "points": 5000
	}'

## Get

	curl --location --request GET 'http://localhost:8080/checkBalances'
