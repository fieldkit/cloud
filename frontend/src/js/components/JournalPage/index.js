
import ControlPanelContainer from '../../containers/common/ControlPanelContainer'
import Post from './Post'
import React, { PropTypes } from 'react'

class JournalPage extends React.Component {

  shouldComponentUpdate (props) {
    return this.props.documents.size === 0
  }

  render () {
    const {
      documents,
      currentExpeditionID,
      updateDate
    } = this.props

    return (
      <div className="journal-page page">
        <div
          className="journal-page_map-overlay"
        />
        <ControlPanelContainer/>
        <div
          className="journal-page_content"
        >
          <ul className="journal-page_content_posts">
            {
              documents
                .sortBy(d =>  d.get('date'))
                .reverse()
                .map(d => {
                  return (
                    <li
                      className="journal-page_content_posts_post"
                      key={ 'post-' + d.get('id') }
                    >
                      <Post
                        data={ d }
                        currentExpeditionID={ currentExpeditionID }
                        updateDate={ updateDate }
                      />
                    </li>
                  )
                })
            }
          </ul>
        </div>
      </div>
    )
  }

}

JournalPage.propTypes = {

}

export default JournalPage