import React, {PropTypes, Component} from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

export default class WebGLOverlay extends Component {

  render () {
    const {
      focusedDocument,
      openLightbox
    } = this.props

    return (
      <div
        className="dom-overlay"
        style={{
          cursor: !!focusedDocument ? 'pointer' : 'auto'
        }}
        onClick={ () => { if (!!focusedDocument) openLightbox(focusedDocument.get('id')) }}
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
                <div className="dom-overlay_popup_content_icon"></div>
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
