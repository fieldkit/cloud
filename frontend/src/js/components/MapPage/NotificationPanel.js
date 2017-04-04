import React, { PropTypes } from 'react'
import { dateToString } from '../../utils.js'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import { is } from 'immutable'

import sensorIcon from '../../../img/icon-sensor-red.svg'

class NotificationPanel extends React.Component {
  shouldComponentUpdate (props) {
    return !is(this.props.currentDocuments, props.currentDocuments)
  }

  render () {
    const {
      currentDocuments
    } = this.props
    
    return (
      <div className="notification-panel">
        <ReactCSSTransitionGroup
          transitionName="transition"
          transitionEnter={true}
          transitionEnterTimeout={500}
          transitionLeave={true}
          transitionLeaveTimeout={500}
        >
          {
            currentDocuments
              .map(d => {
                return (
                  <div class="notification-panel_post">
                    <div className="notification-panel_post_content">
                      <div className="notification-panel_post_content_icon">
                        <img src={ '/' + sensorIcon } width="100%"/>
                      </div>
                      <div>{ d.get('id') }</div>
                      <div>{ d.get('date') }</div>
                    </div>
                  </div>
                )
              })
          }
        </ReactCSSTransitionGroup>
      </div>
    )
  }

}

NotificationPanel.propTypes = {

}

export default NotificationPanel
