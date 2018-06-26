Simple utility to add current temperature and humidity to an existing SmartThings Influx database so that I can overlay
outside info with the inside info. This utilizes the ST database structure that is 
outlined [here](http://codersaur.com/2016/04/smartthings-data-visualisation-using-influxdb-and-grafana/)

Uses:

https://github.com/influxdata/influxdb/tree/master/client

https://github.com/briandowns/openweathermap

You will need an Open Weathermap API key.

Create a config file in `~/.weathertoinflux` called `weathertoinflux.yml` and make it look like:

    #
    # Make the cities fields strings (ie. put them in quotes). It's too variable to parse otherwise.
    #
    openweather:
      apiKey: YOUR_KEY_HERE
      cities:
        - type: "id"
          id: "5039080"
          label: "Outside"
        - type: "latlon"
          lat: "10.1735"
          lon: "-84.733"
          label: "San Buenas"
        - type: "zip"
          zip: "50533"
          countryCode: "US"
          label: "Eagle Grove"
    
    db:
      address: http://localhost:8086
      name: SmartThings
      username: username
      password: password
    
Obviously change things as needed. The cities array can be any mix of things, just make them all strings, I'm converting
in the code. It's just easier to do it that way and I don't want to complicate this any more than is needed for a 
quick and dirty utility.

This puts the data into a InfluxDB database that also aggregates some SmartThings data, therefore the fields and such
are already defined for the thermostats. In particular:

* Series name: temperature
  * tagKey
    * deviceId
    * deviceName
    * groupId
    * groupName
    * hubId
    * hubName
    * locationId
    * locationName
    * unit

  * fieldKey 
    * value: float

* Series name: humidity
  * tagKey
    * deviceId
    * deviceName
    * groupId
    * groupName
    * hubId
    * hubName
    * locationId
    * locationName
    * unit

  * fieldKey
    * value: float



	 