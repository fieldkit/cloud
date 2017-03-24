import React, { PropTypes } from 'react'
import TimelineContainer from '../../containers/common/TimelineContainer'
import HeaderContainer from '../../containers/Root/HeaderContainer'
import MapContainer from '../../containers/Map'
import ExpeditionPanel from './ExpeditionPanel'
import LightboxContainer from '../../containers/common/Lightbox'


export default class Root extends React.Component {
  constructor (props) {
    super(props)
    this.state = {}
  }

  render () {
    const {
      expeditionFetching,
      documentsFetching,
      currentExpedition,
      expeditionPanelOpen,
      documents,
      project,
      expeditions,
      closeExpeditionPanel
    } = this.props

    const currentPage = location.pathname.split('/').filter(p => !!p && p !== currentExpedition)[0] || 'map'

    return (
      <div
        className="root"
        onMouseMove={(e) => {
          const x = e.nativeEvent.clientX
          const y = e.nativeEvent.clientY
          if (currentPage === 'map') {
            this.props.setMousePosition(x, y)
          }
        }}
        onMouseOut={ () => {
          if (currentPage === 'map') {
            this.props.setMousePosition(-1, -1)
          }
        }}
      >
        { !expeditionFetching && !documentsFetching &&
          <div>
            <ExpeditionPanel
              active={ expeditionPanelOpen }
              closeExpeditionPanel={ closeExpeditionPanel }
              expeditions={ expeditions }
              project={ project }
              currentExpedition={ currentExpedition }
            />
            <MapContainer/>
            <div className="root_content">
              <HeaderContainer/>
            {/*
              <LightboxContainer/>
            */}
              <TimelineContainer/>
              {
                !documents || documents.size === 0 &&
                <div className="root_no-document">
                  This expedition doesn't seem to have any document yet...
                </div>
              }
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
