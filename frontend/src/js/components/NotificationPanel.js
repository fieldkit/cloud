
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
          transitionName="post-transition"
          transitionEnter={false}
          transitionEnterTimeout={500}
          transitionLeave={false}
          transitionLeaveTimeout={500}
        >
          {
            currentDocuments
              .filter(d => {
                return true
              })
              .sort((d1, d2) => {
                return d2.get('date') > d1.get('date')
              })
              .map(d => {
                return (
                  <div class="notification-panel_post">
                    <div className="notification-panel_post_icon">
                    </div>
                    <div>{d.get('id')}</div>
                    <div>{d.get('date')}</div>
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