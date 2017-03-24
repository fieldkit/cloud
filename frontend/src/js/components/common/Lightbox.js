
import React from 'react'
import { dateToString } from '../../utils.js'
import { Link } from 'react-router'

import iconLocation from '../../../img/icon-location.png'

class Lightbox extends React.Component {
  render () {
    const {
      data,
      currentExpeditionID,
      updateDate
    } = this.props

    return <div></div>

    return (
      <div className="lightbox">
        <div className="lightbox_map-overlay"/>
        <div className="lightbox_content">
          HELLO WORLD LIGHTBOX
          <div className="lightbox_content_type">
          </div>
          <div className="lightbox_content_main">
            <div className="post_main_meta">
              <div className="post_main_meta_data">
                { dateToString(new Date(data.get('date'))) }
              </div>
              <div className="post_main_meta_geo">
                <Link
                  to={ '/' + currentExpeditionID + '/map' }
                  onClick={ () => updateDate(data.get('date'), 'pause') }
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