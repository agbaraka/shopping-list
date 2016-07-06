Shopping list
===========

This is a shopping list manager hosted on Google App Engine composed of two App Engine modules:
 - A frontend written using [Polymer][1]
 - A backend written in [Go][2] 

## Running locally

To run this application locally install the [Go App Engine SDK][3] and then execute:

```
$ goapp serve dispatch.yaml frontend/app.yaml backend/app.yaml
```
[1]: https://www.polymer-project.org
[2]: https://golang.org
[3]: https://cloud.google.com/appengine/downloads