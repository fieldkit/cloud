// @flow

type LocationType = {
  latitude?: number,
  longitude?: number,
  altitude?: number,
}

export class Location {
  latitude: number;
  longitude: number;
  altitude: number;
  constructor(obj: LocationType){
    this.latitude = obj.latitude ?
      this.latitude = obj.latitude : 0
    this.longitude = obj.longitude ?
      this.longitude = obj.longitude : 0
    this.altitude = obj.altitude ?
      this.altitude = obj.altitude : 0
  }

  describe(): {[string]: string} {
    return {
        "latitude": "number",
        "longitude": "number",
        "altitude": "number",
    }
  }

  describeAll(ns: string = ""): {[string]: string} {
    let fields = {}
    let new_ns = ""
      if(this.latitude){
        if(ns.length > 0){
          fields[ns+".latitude"] = "number"
        } else {
          fields["latitude"] = "number"
        }
        
      }
      if(this.longitude){
        if(ns.length > 0){
          fields[ns+".longitude"] = "number"
        } else {
          fields["longitude"] = "number"
        }
        
      }
      if(this.altitude){
        if(ns.length > 0){
          fields[ns+".altitude"] = "number"
        } else {
          fields["altitude"] = "number"
        }
        
      }
    return fields
  } 
}

type WeatherStationType = {
  temperature?: number,
  humidity?: number,
  pressure?: number,
}

export class WeatherStation {
  temperature: number;
  humidity: number;
  pressure: number;
  constructor(obj: WeatherStationType){
    this.temperature = obj.temperature ?
      this.temperature = obj.temperature : 0
    this.humidity = obj.humidity ?
      this.humidity = obj.humidity : 0
    this.pressure = obj.pressure ?
      this.pressure = obj.pressure : 0
  }

  describe(): {[string]: string} {
    return {
        "temperature": "number",
        "humidity": "number",
        "pressure": "number",
    }
  }

  describeAll(ns: string = ""): {[string]: string} {
    let fields = {}
    let new_ns = ""
      if(this.temperature){
        if(ns.length > 0){
          fields[ns+".temperature"] = "number"
        } else {
          fields["temperature"] = "number"
        }
        
      }
      if(this.humidity){
        if(ns.length > 0){
          fields[ns+".humidity"] = "number"
        } else {
          fields["humidity"] = "number"
        }
        
      }
      if(this.pressure){
        if(ns.length > 0){
          fields[ns+".pressure"] = "number"
        } else {
          fields["pressure"] = "number"
        }
        
      }
    return fields
  } 
}

type FieldkitType = {
  timestamp?: number,
  cell_voltage?: number,
  state_of_charge?: number,
  location?: LocationType,
  weather_station?: WeatherStationType,
}

export class Fieldkit {
  timestamp: number;
  cell_voltage: number;
  state_of_charge: number;
  location: ?Location;
  weather_station: ?WeatherStation;
  constructor(obj: FieldkitType){
    this.timestamp = obj.timestamp ?
      this.timestamp = obj.timestamp : 0
    this.cell_voltage = obj.cell_voltage ?
      this.cell_voltage = obj.cell_voltage : 0
    this.state_of_charge = obj.state_of_charge ?
      this.state_of_charge = obj.state_of_charge : 0
    this.location = obj.location ?
      new Location(obj.location) : null
    this.weather_station = obj.weather_station ?
      new WeatherStation(obj.weather_station) : null
  }

  describe(): {[string]: string} {
    return {
        "timestamp": "number",
        "cell_voltage": "number",
        "state_of_charge": "number",
        "location": "Location",
        "weather_station": "WeatherStation",
    }
  }

  describeAll(ns: string = ""): {[string]: string} {
    let fields = {}
    let new_ns = ""
      if(this.timestamp){
        if(ns.length > 0){
          fields[ns+".timestamp"] = "number"
        } else {
          fields["timestamp"] = "number"
        }
        
      }
      if(this.cell_voltage){
        if(ns.length > 0){
          fields[ns+".cell_voltage"] = "number"
        } else {
          fields["cell_voltage"] = "number"
        }
        
      }
      if(this.state_of_charge){
        if(ns.length > 0){
          fields[ns+".state_of_charge"] = "number"
        } else {
          fields["state_of_charge"] = "number"
        }
        
      }
      if(this.location){
        if(ns.length > 0){
          fields[ns+".location"] = "Location"
        } else {
          fields["location"] = "Location"
        }
        
          if(ns.length > 0){ new_ns = ns + ".location" } else { new_ns = "location" } 
          Object.assign(fields,this.location.describeAll(new_ns))
        
      }
      if(this.weather_station){
        if(ns.length > 0){
          fields[ns+".weather_station"] = "WeatherStation"
        } else {
          fields["weather_station"] = "WeatherStation"
        }
        
          if(ns.length > 0){ new_ns = ns + ".weather_station" } else { new_ns = "weather_station" } 
          Object.assign(fields,this.weather_station.describeAll(new_ns))
        
      }
    return fields
  } 
}
