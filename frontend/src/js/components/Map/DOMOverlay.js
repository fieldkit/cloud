import React, {PropTypes, Component} from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import sensorIcon from '../../../img/icon-sensor-red.svg'

export default class WebGLOverlay extends Component {
  constructor (props) {
    super(props)
    this.state = {}
    this.onClick = this.onClick.bind(this)
  }

  onClick () {
    const {
      focusedDocument,
      openLightbox
    } = this.props
    if (!!focusedDocument) openLightbox(focusedDocument.get('id'))
  }

  render () {
    const {
      focusedDocument
    } = this.props

    return (
      <div
        className="dom-overlay"
        style={{
          cursor: !!focusedDocument ? 'pointer' : 'auto'
        }}
        onClick={ this.onClick }
      >
        <ReactCSSTransitionGroup
          transitionName="transition"
          transitionEnter={true}
          transitionEnterTimeout={200}
          transitionLeave={true}
          transitionLeaveTimeout={200}
        >
        {
          !!focusedDocument &&
            <div 
              class="dom-overlay_popup"
              style={{
                left: focusedDocument.get('x'),
                top: focusedDocument.get('y')
              }}
            >
              <div className="dom-overlay_popup_content">
                <div className="dom-overlay_popup_content_icon">
                  <img src={ '/' + sensorIcon } width="100%"/>
                </div>
                <div>{focusedDocument.get('id')}</div>
                <div>{focusedDocument.get('date')}</div>
              </div>
            </div>
        }
        </ReactCSSTransitionGroup>
      </div>
    )
  }

}

WebGLOverlay.propTypes = {
}
