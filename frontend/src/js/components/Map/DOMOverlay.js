import React, { PropTypes, Component } from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import sensorIcon from '../../../img/icon-sensor-red.svg'
import twitterIcon from '../../../img/icon-twitter.png'

export default class WebGLOverlay extends Component {
    constructor(props) {
        super(props)
        this.state = {}
        this.onClick = this.onClick.bind(this)
    }

    onClick() {
        const {focusedDocument, openLightbox} = this.props
        if (!!focusedDocument) openLightbox(focusedDocument.get('id'))
    }

    render() {
        const {focusedDocument} = this.props

        let title,
            body,
            icon,
            extra;
        let d = focusedDocument;
        if (focusedDocument) {
            if (d.get("user")) {
                title = `@${d.get("user").get("screen_name")}`
                body = d.get("text")
                icon = `/${twitterIcon}`
            } else if (d.get("GPSSpeed")) {
                title = `Sensor (${d.get("SampleType")})`
                body = `GPS Speed: ${d.get("GPSSpeed")}`
                icon = `/${sensorIcon}`
            } else {
                title = "Conservify Sensor"
                body = "Humidity: " + d.get("hum")
                icon = `/${sensorIcon}`
            }
        } else {
            title = ""
            body = ""
            icon = ""
        }
        return (
            <div className="dom-overlay" style={ { cursor: !!focusedDocument ? 'pointer' : 'auto' } } onClick={ this.onClick }>
                <ReactCSSTransitionGroup transitionName="transition" transitionEnter={ true } transitionEnterTimeout={ 200 } transitionLeave={ true }
                    transitionLeaveTimeout={ 200 }>
                    { !!focusedDocument &&
                      <div class="dom-overlay_popup" style={ { left: focusedDocument.get('x'), top: focusedDocument.get('y') } }>
                          <div className="dom-overlay_popup_content">
                              <div className="dom-overlay_popup_content_icon">
                                  <img src={ icon } width="100%" />
                              </div>
                              <div className="dom-overlay_popup_title">
                                  { title }
                              </div>
                              <div>
                                  { body }
                              </div>
                          </div>
                      </div> }
                </ReactCSSTransitionGroup>
            </div>
        )
    }

}

WebGLOverlay.propTypes = {
}
