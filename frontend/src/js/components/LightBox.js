
import React, { PropTypes } from 'react'
import Post from './Post'

const LightBox = (props) => {
  const { active, post, closeLightBox, show360Picture } = props
  var height = {height: window.innerHeight - 200}

  if (active && post) {
    return (
      <div id="lightBox">
        <div id="feed" style={height}>
          <Post format="full" data={post} closeLightBox={closeLightBox} show360Picture={show360Picture}>
            {post.content}
          </Post>
        </div>
      </div>
    )
  } else {
    return null
  }
}

LightBox.propTypes = {
  active: PropTypes.bool.isRequired,
  post: PropTypes.object,
  closeLightBox: PropTypes.func.isRequired,
  show360Picture: PropTypes.func.isRequired
}
export default LightBox

