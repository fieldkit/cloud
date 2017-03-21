
import ControlPanelContainer from '../../containers/common/ControlPanelContainer'
import Post from './Post'
import React, { PropTypes } from 'react'
import I from 'immutable'
import { map } from '../../utils.js'

class JournalPage extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      postDimensions: I.Map(),
      postMarginBottom: 80,
      contentPadding: 0
    }
    this.onScroll = this.onScroll.bind(this)
    this.updateScrollTop = this.updateScrollTop.bind(this)
    this.updatePostLayout = this.updatePostLayout.bind(this)
    this.updateDate = this.updateDate.bind(this)
    this.shouldThisUpdate = this.shouldThisUpdate.bind(this)
  }

  shouldThisUpdate (props, state, prev) {
    const docsHaveUpdated = props.documents.map(d => d.get('id')).join('') !== this.props.documents.map(d => d.get('id')).join('')
    const docDimensionsHaveUpdated = state.postDimensions !== this.state.postDimensions
    const forcedDateUpdated = props.currentDate !== this.props.currentDate && (prev ? this.props.forceDateUpdate : props.forceDateUpdate)
    return docsHaveUpdated || docDimensionsHaveUpdated || forcedDateUpdated
  }

  updateDate () {
    const {
      postDimensions
    } = this.state
    const {
      documents,
      updateDate
    } = this.props
    const scrollTop = this.refs.content.scrollTop
    let firstPost = null
    let secondPost = null
    let scrollRatio = 0
    postDimensions
      .sortBy(p => p.get(0))
      .forEach((p, i) => {
        if (!!firstPost && !!secondPost) return false
        if (!!firstPost) {
          secondPost = documents.find(d => d.get('id') === i.slice(5))
          return true
        }
        if (p.get(0) <= scrollTop && p.get(0) + p.get(1) > scrollTop) {
          firstPost = documents.find(d => d.get('id') === i.slice(5))
          scrollRatio = map(scrollTop, p.get(0), p.get(0) + p.get(1), 0, 1)
          return true
        }
      })
    const currentDate = map(scrollRatio, 0, 1, firstPost.get('date'), secondPost.get('date'))
    updateDate(currentDate)
  }

  updateScrollTop (prevProps, prevState) {
    const {
      documents,
      currentDate
    } = this.props
    const {
      postDimensions
    } = this.state

    if (this.shouldThisUpdate(prevProps, prevState, true)) {
      const timeFrame = [
        documents.last().get('date'),
        documents.first().get('date')
      ]
      let firstDocument = null
      let secondDocument = null
      documents.reverse().forEach((d, i) => {
        if (d.get('date') > currentDate) {
          secondDocument = d
          return false
        }
        firstDocument = d
      })
      const timeRatio = map(currentDate, firstDocument.get('date'), secondDocument.get('date'), 0, 1)
      this.refs.content.scrollTop = map(timeRatio, 0, 1, postDimensions.getIn(['post-' + firstDocument.get('id'), 0]), postDimensions.getIn(['post-' + secondDocument.get('id'), 0]))
    }
  }

  updatePostLayout (prevProps) {
    const {
      documents
    } = this.props

    if ((!prevProps.documents || prevProps.documents.size === 0) && !!documents && documents.size > 0) {
      const postDimensions = {}
      documents.forEach(d => {
        const id = 'post-' + d.get('id')
        postDimensions[id] = [this.refs[id].refs.container.offsetTop, this.refs[id].refs.container.offsetHeight + this.state.postMarginBottom]
      })
      const lastPost = Array.prototype.slice.call(this.refs.posts.childNodes).pop()
      const contentPadding = this.refs.content.offsetHeight - lastPost.offsetHeight - this.state.postMarginBottom - 1

      this.setState({
        ...this.state,
        postDimensions: I.fromJS(postDimensions),
        contentPadding
      })
    }
  }

  shouldComponentUpdate (nextProps, nextState) {
    return this.shouldThisUpdate(nextProps, nextState)
  }

  componentDidUpdate (prevProps, prevState) {
    this.updateScrollTop(prevProps, prevState)
    this.updatePostLayout(prevProps, prevState)
  }

  onScroll (e) {
    this.updateDate()
  }

  render () {
    const {
      documents,
      currentExpeditionID,
      currentDate,
      updateDate
    } = this.props

    const {
      postMarginBottom,
      contentPadding
    } = this.state

    return (
      <div className="journal-page page">
        <div
          className="journal-page_map-overlay"
        />
        <ControlPanelContainer/>
        <div
          className="journal-page_content"
          onScroll={ this.onScroll }
          ref="content"
        >
          <ul
            className="journal-page_content_posts"
            ref="posts"
          >
            {
              documents
                .map(d => {
                  return (
                    <li
                      className="journal-page_content_posts_post"
                      key={ 'post-' + d.get('id') }
                      style={{
                        marginBottom: postMarginBottom
                      }}
                    >
                      <Post
                        data={ d }
                        currentExpeditionID={ currentExpeditionID }
                        updateDate={ updateDate }
                        ref={ 'post-' + d.get('id') }
                      />
                    </li>
                  )
                })
            }
          </ul>
          <div
            style={{
              height: contentPadding
            }}
          />
        </div>
      </div>
    )
  }

}

JournalPage.propTypes = {

}

export default JournalPage