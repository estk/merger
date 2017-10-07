# merger
## Synopsis
The merger project is meant to be a demonstration that it is simple to buffer tracking events and merge them when an impression is fired. An example of this would be in a micro service architecture where each service attaches some domain information to a trackable item, eventually if the user interacts with said trackable item this service will merge all of the data together and potentially push to a data-lake.

## Architecture
This project leverages grpc communication in order to accept partial events and complete events (impressions). Of course it could also consume from a kafka stream of events.

## Status
This is a proof of concept, no error handling has been done, and obviously this should be backed by an actual datastore.
