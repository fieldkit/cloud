import React, { PropTypes } from 'react'
import TimelineContainer from '../containers/TimelineContainer'
import HeaderContainer from '../containers/HeaderContainer'
import MapContainer from '../containers/MapContainer'

export default class Root extends React.Component {
  render () {

    const { expeditionFetching, documentsFetching } = this.props

    return (
      <div className="root">
        { !expeditionFetching && !documentsFetching &&
          <div>
            <MapContainer/>
            <div className="root_content">
              <HeaderContainer/>
              <TimelineContainer/>
              { this.props.children }
            </div>
          </div>
        }
        { expeditionFetching &&
          <div
            style={{
              color: 'black'
            }}
          >
            fetching expedition...
          </div>
        }
        { documentsFetching &&
          <div
            style={{
              color: 'black'
            }}
          >
            fetching documents...
          </div>
        }
      </div>
    )
  }
}

Root.propTypes = {}
