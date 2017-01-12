import React, { PropTypes } from 'react'
import MapGL from 'react-map-gl'

const MapBox = () => {

  const MAPBOX_ACCESS_TOKEN = 'pk.eyJ1IjoiaWFhYWFuIiwiYSI6ImNpbXF1ZW4xOTAwbnl3Ymx1Y2J6Mm5xOHYifQ.6wlNzSdcTlonLBH-xcmUdQ';
  const MAPBOX_STYLE = 'mapbox://styles/iaaaan/cioxda1tz000cbnm13jtwhl8q';

  return (
    <div id="mapbox">
      <MapGL 
        mapStyle={MAPBOX_STYLE}
        mapboxApiAccessToken={MAPBOX_ACCESS_TOKEN}
        width={window.innerWidth} height={window.innerHeight} latitude={37.7577} longitude={-122.4376}
        zoom={8} onChangeViewport={(viewport) => {
          const {latitude, longitude, zoom} = viewport;
        }}
      />
    </div>
  )
}

Okavango.propTypes = {}
export default Map

