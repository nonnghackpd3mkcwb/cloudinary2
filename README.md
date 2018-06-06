
# Go Thumbnail Service + Tests

## Try first (on Heroku)
from your browser go to heroku app [here](https://pure-earth-19327.herokuapp.com/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=1024&height=400)

![img](assets/lena_hi.png)


## Local Tests
install the one dependecy in this project with 
```
go get -u github.com/nfnt/resize
```

than run 
```
go test -v
```

This will run the following tests:
 - test "good flow" -  Go over all kinds of widths+heights combinations and make sure result is 200
 - test non jpeg queries
 - test missing/bad format paramteres in the url

find those on `api_test.go`

### If this was going to production I would add tests that actually check the content of the images being returned by the service for different resolutions. 


## Running Heroku Locally

```sh
$ go get -u https://github.com/nonnghackpd3mkcwb/cloudinary2
$ cd $GOPATH/src/github.com/nonnghackpd3mkcwb/cloudinary2
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

You should also install [govendor](https://github.com/kardianos/govendor) if you are going to add any dependencies to the sample app.

## Question:
###Your thumbnail service is getting traction and CNN wants to use it for
###their website (which contains thousands of new original images every day).
###How would you expand the service, it’s infrastructure and/or architecture to handle the
###large scale?

There are two Microservices which have different requirements and need to scale out differently. 
So I would separate the code I created to:
- Microservice 1 - A first layer of a Load balancer + Stateless Http server containers to accept and hold the connections and finally return the response. Those container will send image processing requests to the second layer for processing. 
- A Queuing service (rabbitmq?) for connectivity between the microservices
- Microservice 2 - A second layer of containers that accept those image processing requests and do the actual processing. 

We can do some optimization like caching the images to popular dimentions (iphone, android, pc) and avoid doing the actual processing for each and every request. 
