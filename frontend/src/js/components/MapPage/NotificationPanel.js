import React, { PropTypes } from 'react'
import { dateToString } from '../../utils.js'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import { is } from 'immutable'

import sensorIcon from '../../../img/icon-sensor-red.svg'
import twitterIcon from '../../../img/icon-twitter.png'

class NotificationPanel extends React.Component {
  shouldComponentUpdate (props) {
    return !is(this.props.currentDocuments, props.currentDocuments)
  }

  render () {
    const {
      currentDocuments
    } = this.props

    let panels = currentDocuments
                  .map(d => { 
                      let title, body, icon, extra;
                      if(d.get("user")){
                        title = `@${d.get("user").get("screen_name")}` 
                        body = d.get("text")
                        icon = `/${twitterIcon}`
                        if(d.getIn(["entities","media"])){
                            let url = d.getIn(["entities","media",0,"media_url_https"])
                            extra = (<img className='notification-panel_twitter-media' src={url}/>)
                        } else {
                            extra = (<span></span>)
                        }
                      } else if (d.get("GPSSpeed")) {
                        title = `Sensor (${d.get("SampleType")})`
                        body = `GPS Speed: ${d.get("GPSSpeed")}`
                        icon = `/${sensorIcon}`
                        extra = (<span></span>)
                      } else {
                        title = `Conservify Sensor`
                        body = "Humidity: " + d.get("hum")
                        icon = `/${sensorIcon}`
                        extra = (<span></span>)
                      }
                      return (
                          <div class="notification-panel_post">
                            <div className="notification-panel_post_content">
                              <div className="notification-panel_post_content_icon">
                                <img src={ icon } width="100%"/>
                              </div>
                              <div className="notification-panel_title">{ title }</div>
                              <div>{ body }</div>
                              {extra}
                            </div>
                          </div>
                      )
                  })
    
    return (
      <div className="notification-panel">
        <ReactCSSTransitionGroup
          transitionName="transition"
          transitionEnter={true}
          transitionEnterTimeout={500}
          transitionLeave={true}
          transitionLeaveTimeout={500}
        >
          {panels}
        </ReactCSSTransitionGroup>
      </div>
    )
  }

}

NotificationPanel.propTypes = {

}

export default NotificationPanel
