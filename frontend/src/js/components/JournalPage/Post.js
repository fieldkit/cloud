
import React from 'react'
import { dateToString } from '../../utils.js'
import { Link } from 'react-router'
import { is } from 'immutable'

import iconLocation from '../../../img/icon-location.png'

class Post extends React.Component {

  shouldComponentUpdate (props) {
    return !is(this.props.data, props.data)
  }

  render () {
    const {
      data,
      currentExpeditionID,
      updateDate
    } = this.props

    return (
      <div className="post" ref="container">
        <div className="post_type">
          Tweet
        </div>
        <div className="post_main">
          <div className="post_main_content">
            { 'This is a post: ' + data.get('id') } 
          </div>
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
        <div className="post_actions">
        </div>
      </div>
    )
  }
}

export default Post