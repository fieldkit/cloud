
import React, {PropTypes} from 'react'
import Feed from './Feed'
import ControlPanelContainer from '../containers/ControlPanelContainer.js'
// import autobind from 'autobind-decorator'

class JournalPage extends React.Component {
  
  render () {
    const {posts, checkFeedContent, expedition} = this.props
    return (
      <div className='page' id="journalPage">
        <ControlPanelContainer/>
        <Feed posts={posts} checkFeedContent={checkFeedContent} expedition={expedition}/>
      </div>
    )
  }
}

JournalPage.propTypes = {
  posts: PropTypes.array.isRequired,
  expedition: PropTypes.object,
  checkFeedContent: PropTypes.func.isRequired
}

export default JournalPage
