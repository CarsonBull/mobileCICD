# mobileCICD
Mobile CI CD system built in go


### Webapp

Simple webapp to allow users to manage their mobile CI/CD pipeline

##### Roadmap

* Add a way for users to define their pipelines
* Remove storage and logic around the pipeline from the webapp. 
* Allow the webapp to connect to different pipelines? (Still not entirely sure wht this would look like or what it means)

### Scheduler

Mobile CI/CD scheduler. Offers a RESTful API that allows users to setup a mobile pipeline. 

##### Roadmap

* Use kafka to allow the pipeline to have pubsub events setup
* Allocate Mac machines for use
* Allocate docker containers for custom pipelines
* Plugins for generic build steps
