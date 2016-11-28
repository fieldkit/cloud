
import React, {PropTypes} from 'react'
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'
import autobind from 'autobind-decorator'

class IntroductionBox extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      complete: false,
      contentEnabled: false,
      startDate: 0,
      currentPosts: [],
      posts: [
        {
          content: <p key="0">The Okavango Delta is one of the world’s last great wetland wildernesses.</p>,
          timeRange: [0, 6000]
        },
        {
          content: <p key="1">A team of Ba’Yei, scientists, engineers and adventurers is journeying a 345 kilometers crossing the delta, finding new species and exploring new ground.</p>,
          timeRange: [6000, 13000]
        },
        {
          content: <p key="2">Join us in real-time as we explore<br/><span>the beating heart of our planet.</span></p>,
          timeRange: [13000, 21000]
        }
      ]
    }
  }

  @autobind
  skip () {
    const { enableContent } = this.props
    this.state.complete = true
    this.state.contentEnabled = true
    enableContent()
  }

  @autobind
  tick () {
    const { enableContent } = this.props
    const { posts, startDate } = this.state

    if (location.pathname.indexOf('/map') === -1) {
      this.state.complete = true
      return
    }

    const now = Date.now() - startDate
    var currentPosts = []
    posts.forEach(p => {
      if (p.timeRange[0] <= now && p.timeRange[1] > now) {
        currentPosts.push(p)
      }
    })

    if (now > posts[posts.length - 1].timeRange[1] - 6000 && !this.state.contentEnabled) {
      this.state.contentEnabled = true
      enableContent()
    }

    var flag = true
    if (currentPosts.length !== this.state.currentPosts.length) flag = false
    else {
      for (var i = 0; i < Math.max(currentPosts.length, this.state.currentPosts.length); i++) {
        if (currentPosts[i] !== this.state.currentPosts[i]) {
          flag = false
          break
        }
      }
    }

    if (!flag) {
      this.setState({
        ...this.state,
        currentPosts: currentPosts
      })
    }
    if (now > posts[posts.length - 1].timeRange[1]) {
      this.state.complete = true
      return
    }
    requestAnimationFrame(this.tick)
  }

  componentDidMount () {
    this.setState({
      ...this.state,
      startDate: Date.now()
    })
    // console.log('aga componentDidMount')
    requestAnimationFrame(this.tick)
  }

  render () {
    const { currentPosts, complete } = this.state
    const { animate } = this.props
    const posts = currentPosts.map(p => {
      return p.content
    })

    const container = () => {
      return (
        <ReactCSSTransitionGroup transitionName="notif" transitionEnterTimeout={500} transitionLeaveTimeout={200}>
          {posts}
        </ReactCSSTransitionGroup>
      )
    }

    if (complete || !animate) return <div></div>
    return (
      <div id="IntroductionBox">
        {container()}
        <div id="skip-button" onClick={this.skip}>SKIP</div>
      </div>
    )
  }
}

IntroductionBox.propTypes = {
  enableContent: PropTypes.func.isRequired,
  animate: PropTypes.bool.isRequired
}

export default IntroductionBox
