# Welcome to Mobile CD/CD webapp

Webapp to support building mobile application pipelines in a simple and easy to use manor.


### Start the web server:

####Development 

   Setup your .env file with your github id, secret and client secret    
   Run `./run-script.sh`
   
##### Go to http://localhost:9000/ and you'll see:
   
    "It works"
   
#### Production
   
   // TODO Create docker image that build the app

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        models/       Models go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## How to develop

The easiest way to contribute is to find TODOs and add the needed change.

## Help

* The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/examples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).

