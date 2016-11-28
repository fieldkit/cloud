
import React, {PropTypes} from 'react'
import ReactDOM from 'react-dom'
import Notification from './Notification'

class NotificationPanel extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      notifications: [],
      currentNotifications: [],
      posY: 0
    }
  }

  componentWillReceiveProps (nextProps) {
    const { posts, currentDate } = nextProps
    if (!posts) return

    let start
    let end
    if (currentDate) {
      start = new Date(currentDate.getTime() - 0.25 * (1000 * 3600))
      end = new Date(currentDate.getTime() + 0.25 * (1000 * 3600))
    }

    const notifications = posts
      .filter(post => {
        return post.type !== 'image' || post.properties.Make !== 'RICOH'
      })
      .sort((a, b) => {
        return new Date(a.properties.DateTime).getTime() - new Date(b.properties.DateTime).getTime()
      })

    const currentNotifications = notifications
      .filter(post => {
        const d = new Date(post.properties.DateTime)
        return d.getTime() >= start && d.getTime() < end
      })

    this.setState({
      ...this.state,
      notifications,
      currentNotifications
    })
  }

  shouldComponentUpdate (nextProps, nextState) {
    const currentProperties = {...this.props, ...this.state}
    const nextProperties = {...nextProps, ...nextState}
    return Object.keys(currentProperties).some(p => {
      if (currentProperties[p] !== nextProperties[p]) {
        return true
      } else {
        return false
      }
    })
  }

  componentDidUpdate () {
    const { notifications, currentNotifications } = this.state
    const { currentDate, playback } = this.props
    if (currentNotifications && currentNotifications.length > 0) {
      currentNotifications
        .filter((notification, i) => {
          return i === 0
        })
        .some(notification => {
          return this.setState({
            ...this.state,
            posY: -1 * ReactDOM.findDOMNode(this).querySelector('div.notification.n' + notification.id).offsetTop
          })
        })
    } else {
      const forward = playback === 'fastForward' || playback === 'forward' || playback === 'pause'
      if (forward) {
        for (let i = 0; i < notifications.length; i++) {
          const notification = notifications[i]
          const d = new Date(notification.properties.DateTime)
          if (d.getTime() > currentDate.getTime()) {
            return this.setState({
              ...this.state,
              posY: -1 * ReactDOM.findDOMNode(this).querySelector('div.notification.n' + notification.id).offsetTop
            })
          }
        }
      } else {
        for (let i = notifications.length - 1; i >= 0; i--) {
          const notification = notifications[i]
          const d = new Date(notification.properties.DateTime)
          if (d.getTime() < currentDate.getTime()) {
            return this.setState({
              ...this.state,
              posY: -1 * ReactDOM.findDOMNode(this).querySelector('div.notification.n' + notification.id).offsetTop
            })
          }
        }
      }
    }
  }

  render () {
    const { notifications, currentNotifications } = this.state

    const notificationItems = notifications
      .filter((post, i) => {
        return window.innerWidth > 768 || currentNotifications.indexOf(post) === currentNotifications.length - 1
      })
      .map((post, i) => {
        const active = currentNotifications.indexOf(post) > -1
        switch (post.type) {
          case 'tweet':
            var text = post.properties.Text
            text = text.split(' ').slice(0, text.split(' ').length - 1).join(' ')
            if (post.properties.Tweet) text = post.properties.Tweet.text
            var images = post.properties.Images
              .filter((img, j) => {
                return j === 0
              })
              .map((img, j) => {
                return <img src={active ? img.Url.replace('http://', 'https://') : '#'} width="100%" key={j}/>
              })
            return (
              <Notification
                type={post.type}
                key={post.id}
                active={active}
                id={post.id}
              >
                <p style={{position: window.innerWidth > 768 || post.properties.Images.length === 0 ? 'relative' : 'absolute'}}>{text}</p>
                <div className="images">{images}</div>
              </Notification>
            )
          case 'audio':
            return (
              <Notification
                type={post.type}
                key={post.id}
                id={post.id}
                active={active}
              >
                <div className="title">{post.properties.Title}</div>
              </Notification>
            )
          case 'blog':
            return (
              <Notification
                type={post.type}
                key={post.id}
                id={post.id}
                active={active}
              >
                <div className="title">{post.properties.Title}</div>
              </Notification>
            )
          case 'image':
            var width = 0
            var height = 0
            if (post.properties.Dimensions) {
              width = post.properties.Dimensions[0]
              height = post.properties.Dimensions[1]
            } else if (post.properties.Size) {
              width = post.properties.Size[0]
              height = post.properties.Size[1]
            }
            return (
              <Notification
                type={post.type}
                key={post.id}
                id={post.id}
                active={currentNotifications.indexOf(post) > -1}
              >
                <img className="image" src={active ? post.properties.Url.replace('http://', 'https://') : '#'} width={width} height={height}/>
              </Notification>
            )
        }
      })

    var height = window.innerWidth > 768 ? {height: window.innerHeight - 110} : {}

    return (
      <div id="notificationPanel" style={height}>
        <div
          className="scrollableContainer"
          style={{
            top: this.state.posY
          }}
        >
          {notificationItems}
        </div>
      </div>
    )
  }
}

NotificationPanel.propTypes = {
  posts: PropTypes.array,
  currentDate: PropTypes.object,
  playback: PropTypes.string
}

export default NotificationPanel
