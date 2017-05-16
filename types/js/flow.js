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
}
