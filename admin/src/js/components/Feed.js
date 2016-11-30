
import React, {PropTypes} from 'react'
import Post from './Post'
import * as d3 from 'd3'
import autobind from 'autobind-decorator'

class Feed extends React.Component {
  constructor (props) {
    super(props)
  }

  render () {
    const { posts, checkFeedContent } = this.props

    const format = 'full'
    const postFeed = posts
      .slice(0)
      .sort((a, b) => {
        return b.date - a.date
      })
      .map(post => {
        return (
          <Post format={format} data={post} key={post.key}>
            {post.content}
          </Post>
        )
      })

    return (
      <div id="feed" onWheel={checkFeedContent}>
        {postFeed}
      </div>
    )
  }
}

Feed.propTypes = {
  posts: PropTypes.array.isRequired,
  expedition: PropTypes.object,
  checkFeedContent: PropTypes.func.isRequired
}

export default Feed




