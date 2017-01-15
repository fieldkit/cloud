
import React, { PropTypes } from 'react'
import { dateToString } from '../utils.js'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

class NotificationPanel extends React.Component {

  constructor (props) {
    super(props)
    this.state = {
    }
  }

  render () {

    const { currentDate, selectPlaybackMode, currentDocuments } = this.props

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
              .filter(d => {
                return Math.abs(d.get('date') - currentDate + 100000) < 200000
              })
              .sort((d1, d2) => {
                return d2.get('date') > d1.get('date')
              })
              .map(d => {
                return (
                  <div class="notification-panel_post">
                    <div className="notification-panel_post_content">
                      <div className="notification-panel_post_content_icon"></div>
                      <div>{d.get('id')}</div>
                      <div>{d.get('date')}</div>
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