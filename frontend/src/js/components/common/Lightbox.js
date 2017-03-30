
import React from 'react'
import { dateToString } from '../../utils.js'
import { Link } from 'react-router'

import iconLocation from '../../../img/icon-location.png'

class Lightbox extends React.Component {
  constructor (props) {
    super(props)
    this.state = {}
    this.onClick = this.onClick.bind(this)
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
            { data.get('type') }
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
          </div>
        </div>
      </div>
    )
  }
}

export default Lightbox