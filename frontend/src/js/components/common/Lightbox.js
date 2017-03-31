
import React from 'react'
import { dateToString } from '../../utils.js'
import { Link } from 'react-router'

import iconLocation from '../../../img/icon-location.png'
import iconSensor from '../../../img/icon-sensor-red.png'
import iconBackward from '../../../img/icon-backward-hover.png'
import iconForward from '../../../img/icon-forward-hover.png'

class Lightbox extends React.Component {
  constructor (props) {
    super(props)
    this.state = {}
    this.onClick = this.onClick.bind(this)
    this.openPreviousDocument = this.openPreviousDocument.bind(this)
    this.openNextDocument = this.openNextDocument.bind(this)
  }

  openPreviousDocument (e) {
    const {
      previousDocumentID,
      openLightbox
    } = this.props
    if (!!previousDocumentID) {
      openLightbox(previousDocumentID)
    }
    e.stopPropagation()
  }

  openNextDocument (e) {
    const {
      nextDocumentID,
      openLightbox
    } = this.props
    if (!!nextDocumentID) {
      openLightbox(nextDocumentID)
    }
    e.stopPropagation()
  }

  onClick () {
    const {
      data,
      closeLightbox
    } = this.props
    updateDate(data.get('date'), 'pause')
    closeLightbox()
  }

  onClick () {
    const {
      data,
      closeLightbox
    } = this.props
    updateDate(data.get('date'), 'pause')
    closeLightbox()
  }

  render () {
    const {
      data,
      currentExpeditionID,
      updateDate,
      closeLightbox
    } = this.props

    if (!data) return null

    const type = data.get('type')

    return (
      <div className="lightbox">
        <div className="lightbox_map-overlay"/>
        <div
          className="lightbox_content"
          style={{
            padding: `0 ${ 1.388888889 * 12.5 }vw`
          }}
          onClick={ closeLightbox }
        >
          <div className="lightbox_content_type">
            <img src={ '/' + iconSensor } />
          </div>
          <div className="lightbox_content_main">
            <p>
              { `This is a sensor reading: ${ data.get('date') }`}
            </p>
            <div className="post_main_meta">
              <div className="post_main_meta_data">
                { dateToString(new Date(data.get('date'))) }
              </div>
              <div className="post_main_meta_geo">
                <Link
                  to={ '/' + currentExpeditionID + '/map' }
                  onClick={ this.onClick }
                >
                  <img src={ '/' + iconLocation }/>
                </Link>
              </div>
              <div className="post_main_meta_separator"/>
            </div>
          </div>
          <div className="lightbox_content_actions">
            <div className="button"
              onClick={ this.openPreviousDocument }
            >
              <img src={ '/' + iconBackward } />
            </div>
            <div className="button"
              onClick={ this.openNextDocument }
            >
              <img src={ '/' + iconForward } />
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default Lightbox